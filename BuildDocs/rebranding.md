# Gander Social Rebranding Asset Replacement List

## License Compliance Overview

The BlueSky social-app repository is licensed under the **MIT License**, which allows for modification, distribution, and sublicensing. However, the forking guidelines specifically require:

### Forking Requirements from BlueSky README:
- ✅ **Change all branding** in the repository and UI to clearly differentiate from Bluesky
- ✅ **Change any support links** (feedback, email, terms of service, etc) to your own systems
- ✅ **Replace any analytics or error-collection systems** with your own

## 1. Core Configuration Files

### 1.1 Package Configuration
- **File**: `package.json`
    - **Current**: `"name": "gander.social.app"` *(already updated)*
    - **Bundle ID**: Change from `xyz.blueskyweb.app` to `ca.gandersocial.app`
    - **Version**: Maintain independent versioning
    - **Dependencies**: Review Sentry configuration

### 1.2 App Configuration (if exists)
- **Files to check**: `app.json`, `app.config.js`, `app.config.ts`
    - App name: "Gander Social"
    - Package identifier: `ca.gandersocial.app`
    - Bundle identifier: `ca.gandersocial.app`

### 1.3 Build Configuration
- **File**: `eas.json` *(search needed)*
    - Update build profiles
    - Change bundle identifiers
    - Update store configurations

## 2. Visual Assets & Branding

### 2.1 App Icons
**Location**: `./assets/icons/` (referenced in package.json)
- **App icon** (all sizes: 16x16 to 1024x1024)
- **Adaptive icons** (Android foreground/background)
- **Alternative app icons** (BlueSky+ icons mentioned in messages.js)

### 2.2 Splash Screens
- **Splash screen images** (all device sizes)
- **Launch screen assets**
- **Loading indicators**

### 2.3 In-App Visual Elements
**Location**: `./assets/` directory structure
- **Logo variations** (header, footer, compact)
- **Wordmark/text logo**
- **Favicon** (for web version)
- **Social media preview images**

### 2.4 UI Icons & Graphics
**Location**: `./assets/icons/` (SVG optimization in package.json)
- **Navigation icons**
- **Feature-specific icons**
- **Illustration assets**
- **Background patterns/graphics**

## 3. Text Content & Messaging

### 3.1 Internationalization Files
**Location**: Referenced in package.json (`intl:extract`, `intl:compile`)
- **File**: `messages.js` (contains all UI text)
    - Replace "Bluesky" with "Gander" throughout
    - Update service-specific terminology
    - Modify legal links and support references

### 3.2 Key Text Replacements Needed:
```
"Bluesky" → "Gander" / "Gander Social"
"gndr.app" → "gander.social" 
"Bluesky Social PBC" → "Gander Social"
"Bluesky Social Terms of Service" → "Gander Social Terms of Service"
"@gndr.app" → "@gander.social"
"security@gndr.app" → "security@gander.social"
```

## 4. Web Configuration

### 4.1 Webpack Configuration
**File**: `webpack.config.js`
- **Sentry organization**: Change from `'blueskyweb'` to appropriate Gander org
- **Project name**: Update from `'app'` to Gander-specific project name

### 4.2 Web Assets
- **Web favicon**
- **Web app manifest**
- **PWA icons**
- **OG/meta images**

## 5. Platform-Specific Assets

### 5.1 iOS Assets (./ios/ directory)
- **App icons** (AppIcon.appiconset)
- **Launch images**
- **App Store assets**
- **Bundle identifier** in `Info.plist`
- **App name** in `Info.plist`

### 5.2 Android Assets (./android/ directory)
- **App icons** (mipmap directories)
- **Adaptive icons**
- **Splash screens**
- **Package name** in `build.gradle`
- **App name** in `strings.xml`

## 6. Code References

### 6.1 Bundle/Package Identifiers
**Current references**: `xyz.blueskyweb.app`
**Replace with**: `ca.gandersocial.app`

**Files to check**:
- All build configuration files
- Native platform configurations
- Testing configurations (e2e, performance tests)

### 6.2 Analytics & Error Reporting
- **Sentry configuration**: Replace with Gander-specific Sentry project
- **Analytics IDs**: Replace any BlueSky-specific tracking IDs
- **Crash reporting**: Update service configurations

## 7. Legal & Support Documentation

### 7.1 Terms & Policies
- **Terms of Service** links and content
- **Privacy Policy** links and content
- **Support/Help** documentation
- **Community guidelines**

### 7.2 Contact Information
- **Support email addresses**
- **Feedback mechanisms**
- **Bug report channels**
- **Legal contact information**

## 8. Server & API Configuration

### 8.1 API Endpoints
- **Default server URLs** (if hardcoded)
- **Service discovery endpoints**
- **CDN configurations**

### 8.2 Domain References
- Replace `gndr.app` domains with `gander.social`
- Update any hardcoded BlueSky service URLs

## 9. Store Listings & Marketing

### 9.1 App Store Assets
- **App Store screenshots**
- **App description**
- **App keywords**
- **Developer information**

### 9.2 Marketing Materials
- **Promotional graphics**
- **Feature highlight images**
- **Video previews** (if any)

## 10. Development & Build Tools

### 10.1 Scripts & Automation
**Files**: Various script files referenced in package.json
- Update any script references to BlueSky services
- Modify build and deployment scripts
- Update error collection endpoints

### 10.2 Development Configuration
- **Environment variable defaults**
- **Development server configurations**
- **Debug/testing service endpoints**

## Implementation Priority

### Phase 1: Critical Branding (Week 1)
1. App icons and splash screens
2. App name and bundle identifiers
3. Core text replacements in messages.js
4. Package.json updates

### Phase 2: Service Integration (Week 2)
1. Analytics and error reporting services
2. API endpoint configurations
3. Legal documentation updates
4. Support system integration

### Phase 3: Platform Optimization (Week 3)
1. Store listing preparations
2. Platform-specific asset optimization
3. Marketing material creation
4. Final testing and validation

## Compliance Checklist

- [ ] All "Bluesky" references replaced with "Gander"
- [ ] Bundle identifiers updated to Gander domains
- [ ] Support links point to Gander systems
- [ ] Analytics/error collection uses Gander services
- [ ] Legal documentation reflects Gander ownership
- [ ] App store listings prepared for Gander brand
- [ ] No BlueSky logos or trademarked assets remain
- [ ] MIT license compliance maintained
- [ ] Attribution to original BlueSky project included where appropriate

## Technical Notes

1. **Icon Optimization**: The project uses SVGO for icon optimization (`"icons:optimize": "svgo -f ./assets/icons"`)
2. **Internationalization**: Full i18n support with lingui framework
3. **Multi-platform**: Supports iOS, Android, and Web
4. **Asset Management**: Uses Expo's asset management system
5. **Bundle Analysis**: Webpack bundle analyzer available for optimization

This comprehensive rebranding ensures full compliance with BlueSky's forking guidelines while establishing Gander Social as a distinct platform.