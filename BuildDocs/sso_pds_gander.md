# Enterprise PDS SSO Implementation Guide for Data Sovereignty

## Overview

This guide provides a secure implementation pattern for Single Sign-On (SSO) integration with private Personal Data Servers (PDS) while maintaining complete data sovereignty for government and business clients.

## Architecture Overview

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   Enterprise    │     │   SSO Bridge     │     │  Private PDS    │
│   Identity      │────▶│   Service        │────▶│  (AT Protocol)  │
│   Provider      │     │                  │     │                 │
└─────────────────┘     └──────────────────┘     └─────────────────┘
        │                        │                         │
        └────────────────────────┴─────────────────────────┘
                    All within Canadian Infrastructure
```

## Implementation Components

### 1. SSO Authentication Service

```typescript
// sso-auth-service.ts
import { SAMLStrategy } from 'passport-saml';
import { Strategy as OIDCStrategy } from 'passport-openidconnect';

export class EnterpriseSSOService {
  private strategies: Map<string, passport.Strategy>;
  
  constructor(private config: SSOConfig) {
    this.initializeStrategies();
  }
  
  private initializeStrategies() {
    // SAML 2.0 for government clients
    if (this.config.saml) {
      this.strategies.set('saml', new SAMLStrategy({
        entryPoint: this.config.saml.entryPoint,
        issuer: this.config.saml.issuer,
        cert: this.config.saml.certificate,
        identifierFormat: 'urn:oasis:names:tc:SAML:2.0:nameid-format:persistent',
        validateInResponseTo: true,
        disableRequestedAuthnContext: false
      }, this.verifySAMLUser.bind(this)));
    }
    
    // OIDC for enterprise clients
    if (this.config.oidc) {
      this.strategies.set('oidc', new OIDCStrategy({
        issuer: this.config.oidc.issuer,
        authorizationURL: this.config.oidc.authorizationURL,
        tokenURL: this.config.oidc.tokenURL,
        userInfoURL: this.config.oidc.userInfoURL,
        clientID: this.config.oidc.clientID,
        clientSecret: this.config.oidc.clientSecret,
        callbackURL: this.config.oidc.callbackURL,
        scope: ['openid', 'profile', 'email']
      }, this.verifyOIDCUser.bind(this)));
    }
  }
  
  async authenticateUser(
    protocol: 'saml' | 'oidc', 
    credentials: any
  ): Promise<EnterpriseUser> {
    const strategy = this.strategies.get(protocol);
    if (!strategy) {
      throw new Error(`Unsupported SSO protocol: ${protocol}`);
    }
    
    return new Promise((resolve, reject) => {
      strategy.authenticate(credentials, (err, user) => {
        if (err) reject(err);
        else resolve(this.mapToEnterpriseUser(user));
      });
    });
  }
  
  private mapToEnterpriseUser(ssoUser: any): EnterpriseUser {
    return {
      id: ssoUser.nameID || ssoUser.sub,
      email: ssoUser.email,
      displayName: ssoUser.displayName || ssoUser.name,
      department: ssoUser.department,
      clearanceLevel: ssoUser.clearanceLevel,
      roles: this.extractRoles(ssoUser),
      attributes: ssoUser.attributes || {}
    };
  }
}
```

### 2. PDS Account Mapping Service

```typescript
// pds-account-mapper.ts
export class PDSAccountMapper {
  constructor(
    private organizationDomain: string,
    private pdsEndpoint: string,
    private keyManager: KeyManagementService
  ) {}
  
  async mapSSOUserToPDS(
    enterpriseUser: EnterpriseUser
  ): Promise<PDSAccount> {
    // Check if PDS account already exists
    let pdsAccount = await this.findExistingPDSAccount(enterpriseUser.id);
    
    if (!pdsAccount) {
      // Create new PDS account
      pdsAccount = await this.createPDSAccount(enterpriseUser);
    }
    
    // Update account metadata
    await this.updateAccountMetadata(pdsAccount, enterpriseUser);
    
    return pdsAccount;
  }
  
  private async createPDSAccount(user: EnterpriseUser): Promise<PDSAccount> {
    // Generate organization-controlled DID
    const did = this.generateOrganizationDID(user.id);
    
    // Create AT Protocol handle
    const handle = this.generateHandle(user);
    
    // Generate signing keys (stored in HSM)
    const keys = await this.keyManager.generateKeys(did);
    
    // Create PDS repository
    const repo = await this.initializePDSRepository(did, keys);
    
    return {
      did,
      handle,
      email: user.email,
      pdsEndpoint: this.pdsEndpoint,
      repository: repo,
      metadata: {
        ssoId: user.id,
        department: user.department,
        clearanceLevel: user.clearanceLevel,
        createdAt: new Date(),
        lastAuthenticated: new Date()
      }
    };
  }
  
  private generateOrganizationDID(userId: string): string {
    // Use did:web for organizational control
    return `did:web:${this.organizationDomain}:users:${userId}`;
  }
  
  private generateHandle(user: EnterpriseUser): string {
    const username = user.email.split('@')[0].toLowerCase();
    return `${username}.${this.organizationDomain}`;
  }
}
```

### 3. Session Management Bridge

```typescript
// session-bridge.ts
export class SessionBridge {
  private sessionStore: SecureSessionStore;
  private readonly SESSION_DURATION = 8 * 60 * 60 * 1000; // 8 hours
  
  constructor(
    private ssoService: EnterpriseSSOService,
    private pdsMapper: PDSAccountMapper,
    private auditLogger: AuditLogger
  ) {
    this.sessionStore = new SecureSessionStore({
      encryption: 'AES-256-GCM',
      storage: 'redis-canadian-region'
    });
  }
  
  async createSession(
    ssoToken: string,
    ipAddress: string,
    userAgent: string
  ): Promise<ATProtoSession> {
    try {
      // 1. Validate SSO token
      const enterpriseUser = await this.ssoService.authenticateUser(
        'saml', 
        { token: ssoToken }
      );
      
      // 2. Map to PDS account
      const pdsAccount = await this.pdsMapper.mapSSOUserToPDS(enterpriseUser);
      
      // 3. Create AT Protocol session
      const session = await this.createATProtoSession(pdsAccount);
      
      // 4. Store session with encryption
      await this.sessionStore.set(session.id, {
        ...session,
        ssoSessionId: ssoToken,
        ipAddress,
        userAgent,
        expiresAt: Date.now() + this.SESSION_DURATION
      });
      
      // 5. Audit log
      await this.auditLogger.log({
        event: 'session_created',
        userId: pdsAccount.did,
        ssoId: enterpriseUser.id,
        ipAddress,
        timestamp: new Date()
      });
      
      return session;
    } catch (error) {
      await this.auditLogger.logError({
        event: 'session_creation_failed',
        error: error.message,
        ipAddress,
        timestamp: new Date()
      });
      throw error;
    }
  }
  
  private async createATProtoSession(
    pdsAccount: PDSAccount
  ): Promise<ATProtoSession> {
    // Generate JWT for AT Protocol
    const accessJwt = await this.generateAccessJWT(pdsAccount);
    const refreshJwt = await this.generateRefreshJWT(pdsAccount);
    
    return {
      id: generateSessionId(),
      did: pdsAccount.did,
      handle: pdsAccount.handle,
      email: pdsAccount.email,
      accessJwt,
      refreshJwt,
      pdsEndpoint: pdsAccount.pdsEndpoint
    };
  }
}
```

### 4. Data Sovereignty Controls

```typescript
// sovereignty-controls.ts
export class DataSovereigntyEnforcer {
  private geoBlocker: GeoBlockingService;
  private encryptionService: EncryptionService;
  
  constructor(private config: SovereigntyConfig) {
    this.geoBlocker = new GeoBlockingService(config.allowedCountries);
    this.encryptionService = new EncryptionService({
      algorithm: 'AES-256-GCM',
      keyStorage: 'canadian-hsm'
    });
  }
  
  async enforceAccessControls(
    request: IncomingRequest,
    session: ATProtoSession
  ): Promise<AccessDecision> {
    // 1. Verify request origin
    const geoCheck = await this.geoBlocker.verifyLocation(request.ip);
    if (!geoCheck.allowed) {
      return { allowed: false, reason: 'geographic_restriction' };
    }
    
    // 2. Verify network requirements
    if (this.config.vpnRequired && !this.isOnVPN(request)) {
      return { allowed: false, reason: 'vpn_required' };
    }
    
    // 3. Verify session validity
    const sessionCheck = await this.verifySession(session);
    if (!sessionCheck.valid) {
      return { allowed: false, reason: 'invalid_session' };
    }
    
    // 4. Apply role-based access
    const rbacCheck = await this.checkRBAC(session, request.resource);
    if (!rbacCheck.allowed) {
      return { allowed: false, reason: 'insufficient_permissions' };
    }
    
    return { allowed: true };
  }
  
  async encryptUserData(data: any, userId: string): Promise<EncryptedData> {
    // Use organization-specific encryption keys
    const encryptionKey = await this.getOrganizationKey(userId);
    
    return this.encryptionService.encrypt(data, encryptionKey, {
      additionalData: { userId, timestamp: Date.now() }
    });
  }
}
```

### 5. Deployment Configuration

```yaml
# enterprise-pds-deployment.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: enterprise-pds-config
  namespace: sovereign-pds
data:
  config.yaml: |
    sso:
      saml:
        enabled: true
        entryPoint: "https://idp.organization.ca/saml/sso"
        issuer: "https://pds.organization.ca"
        certificate: "/certs/saml.crt"
      oidc:
        enabled: false
    
    pds:
      endpoint: "https://pds.organization.ca"
      didMethod: "did:web"
      keyManagement:
        type: "hsm"
        provider: "thales-canadian"
        keyRotationDays: 90
    
    sovereignty:
      dataResidency:
        - region: "canada-central"
          provider: "thinkon"
        - region: "canada-east"
          provider: "thinkon"
      encryption:
        atRest: "AES-256-GCM"
        inTransit: "TLS-1.3"
      networkControls:
        vpnRequired: true
        allowedCountries: ["CA"]
        geoBlocking: true
    
    compliance:
      pipeda: true
      retentionDays: 365
      auditLogging: true
      rightToErasure: true

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: enterprise-pds-sso
  namespace: sovereign-pds
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: sso-bridge
        image: gander/enterprise-pds-sso:latest
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: SOVEREIGN_MODE
          value: "true"
        volumeMounts:
        - name: config
          mountPath: /etc/pds
        - name: certs
          mountPath: /certs
      volumes:
      - name: config
        configMap:
          name: enterprise-pds-config
      - name: certs
        secret:
          secretName: enterprise-certs
```

## Security Considerations

### 1. Key Management

- All signing keys stored in Hardware Security Modules (HSM)
- Regular key rotation (90-day default)
- Backup keys encrypted and stored in separate Canadian facilities
- No keys leave Canadian jurisdiction

### 2. Access Control

- Multi-factor authentication required for all SSO
- Session timeout after 8 hours of inactivity
- IP allowlisting for government networks
- VPN required for remote access

### 3. Audit Trail

```typescript
interface AuditLog {
  timestamp: Date;
  eventType: string;
  userId: string;
  ssoSessionId: string;
  ipAddress: string;
  userAgent: string;
  resource: string;
  action: string;
  result: 'success' | 'failure';
  metadata: Record<string, any>;
}
```

### 4. Emergency Procedures

- Break-glass access for emergencies
- Legal hold capabilities for compliance
- Account suspension without data deletion
- Bulk export for audit/compliance

## Testing Checklist

- [ ] SSO authentication flow works correctly
- [ ] PDS accounts created with proper DIDs
- [ ] Data remains within Canadian borders
- [ ] Encryption at rest and in transit verified
- [ ] Audit logs capture all access
- [ ] Session timeout enforced
- [ ] Geographic restrictions working
- [ ] VPN requirements enforced
- [ ] Key rotation functioning
- [ ] Emergency access procedures tested

## Support and Maintenance

- 24/7 monitoring of authentication services
- Regular security audits (quarterly)
- Compliance reviews (annual)
- Disaster recovery testing (bi-annual)
- Key rotation reminders (automated)

This implementation ensures complete data sovereignty while providing seamless SSO integration for enterprise users.