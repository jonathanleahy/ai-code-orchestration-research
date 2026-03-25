# Localisation Plan

# CRM App Localization Plan

## Priority Languages (Top 5 by Market Size)

1. **Spanish** - 460M speakers globally
2. **French** - 280M speakers globally  
3. **German** - 100M speakers globally
4. **Japanese** - 125M speakers globally
5. **Portuguese** - 260M speakers globally

## What to Translate

### UI Strings
- All user-facing interface elements
- Button labels and tooltips
- Form field placeholders and labels
- Error messages and validation text
- Navigation menu items
- Dashboard widget titles
- Modal dialog content
- Status indicators and badges

### Help Documentation
- User guides and tutorials
- FAQ articles
- Video transcript captions
- Knowledge base articles
- Onboarding walkthroughs

### Legal Documents
- Terms of Service
- Privacy Policy
- Cookie Policy
- EULA (End User License Agreement)
- Data Processing Agreements
- Compliance documentation

### Marketing Materials
- Landing page content
- Email templates
- Social media posts
- Ad copy and banners
- Product descriptions
- Blog content

## What NOT to Translate

### Code Elements
- Variable names and function names
- API endpoint paths
- Database field names
- Configuration keys
- Error codes and status codes
- File names and paths

### Technical Content
- API responses (JSON/XML data)
- Internal logs and debugging output
- Database queries and SQL statements
- Code comments and documentation
- Technical specifications and documentation

### Dynamic Data
- User-generated content
- Customer names and company names
- Product SKUs and identifiers
- Date/time stamps and timestamps
- Numerical values and metrics

## i18n Architecture

### Go Implementation

```go
// Translation structure
type Translation struct {
    ID          string `json:"id"`
    Message     string `json:"message"`
    Description string `json:"description,omitempty"`
}

// Translation manager
type TranslationManager struct {
    translations map[string]map[string]Translation
    defaultLocale string
}

// File structure: locales/{locale}/messages.json
// Example: locales/es/messages.json
{
  "welcome_message": {
    "message": "Bienvenido a nuestra aplicación",
    "description": "Welcome message on login screen"
  },
  "save_button": {
    "message": "Guardar",
    "description": "Save button text"
  }
}

// Usage example
func (tm *TranslationManager) GetTranslation(locale, key string) string {
    if translations, exists := tm.translations[locale]; exists {
        if translation, exists := translations[key]; exists {
            return translation.Message
        }
    }
    return tm.translations[tm.defaultLocale][key].Message
}
```

### TypeScript Implementation

```typescript
// Translation structure
interface Translation {
  id: string;
  message: string;
  description?: string;
}

// Translation manager
class TranslationManager {
  private translations: Record<string, Record<string, Translation>> = {};
  private defaultLocale: string = 'en';
  
  loadTranslations(locale: string, translations: Record<string, Translation>) {
    this.translations[locale] = translations;
  }
  
  t(key: string, locale: string = this.defaultLocale): string {
    return this.translations[locale]?.[key]?.message || 
           this.translations[this.defaultLocale]?.[key]?.message || 
           key;
  }
}

// File structure: src/locales/{locale}/messages.json
// Example: src/locales/es/messages.json
{
  "welcome_message": {
    "message": "Bienvenido a nuestra aplicación",
    "description": "Welcome message on login screen"
  }
}
```

### Translation File Organization

```
locales/
├── en/
│   ├── messages.json
│   ├── legal.json
│   └── help.json
├── es/
│   ├── messages.json
│   ├── legal.json
│   └── help.json
├── fr/
│   ├── messages.json
│   ├── legal.json
│   └── help.json
└── de/
    ├── messages.json
    ├── legal.json
    └── help.json
```

## Date/Currency/Number Formatting per Locale

### Date Formatting
- **Spanish (es)**: dd/MM/yyyy, dd/MM/yyyy HH:mm
- **French (fr)**: dd/MM/yyyy, dd/MM/yyyy HH:mm
- **German (de)**: dd.MM.yyyy, dd.MM.yyyy HH:mm
- **Japanese (ja)**: yyyy/MM/dd, yyyy/MM/dd HH:mm
- **Portuguese (pt)**: dd/MM/yyyy, dd/MM/yyyy HH:mm

### Currency Formatting
- **Spanish**: €1.234,56 (European style)
- **French**: 1 234,56 € (European style)
- **German**: 1.234,56 € (European style)
- **Japanese**: ¥1,234.56 (Japanese yen)
- **Portuguese**: R$ 1.234,56 (Brazilian real)

### Number Formatting
- **Spanish**: 1.234,56 (decimal comma)
- **French**: 1 234,56 (space separator, comma decimal)
- **German**: 1.234,56 (decimal comma)
- **Japanese**: 1,234.56 (comma separator, period decimal)
- **Portuguese**: 1.234,56 (decimal comma)

### Implementation Example (Go)
```go
import "golang.org/x/text/language"
import "golang.org/x/text/message"
import "golang.org/x/text/number"

func FormatCurrency(amount float64, locale string) string {
    tag := language.MustParse(locale)
    p := message.NewPrinter(tag)
    return p.Sprintf("%v", number.Decimal(amount, number.Scale(2)))
}
```

## RTL Support Assessment

### Languages Requiring RTL Support
- **Arabic** (ar) - Right-to-left
- **Hebrew** (he) - Right-to-left
- **Persian** (fa) - Right-to-left
- **Urdu** (ur) - Right-to-left

### Implementation Strategy
1. **CSS Framework**: Use CSS logical properties (start/end instead of left/right)
2. **Layout Direction**: Apply `direction: rtl` to RTL language containers
3. **Component Alignment**: 
   - Text alignment: `text-align: right` for RTL
   - Icon positioning: flip icons in RTL contexts
   - Form layouts: reverse field order
4. **Navigation**: Reverse breadcrumb order, menu positioning

### CSS Example
```css
[dir="rtl"] {
  direction: rtl;
  text-align: right;
}

/* Use logical properties */
.margin-start { margin-inline-start: 1rem; }
.padding-end { padding-inline-end: 1rem; }
```

## Estimated Effort per Language

### Development & Implementation (Initial Setup)
| Language | Hours | Description |
|----------|-------|-------------|
| Spanish | 40 | UI + Legal + Help |
| French | 40 | UI + Legal + Help |
| German | 40 | UI + Legal + Help |
| Japanese | 60 | UI + Legal + Help + Complex formatting |
| Portuguese | 40 | UI + Legal + Help |

### Translation & Quality Assurance
| Language | Hours | Description |
|----------|-------|-------------|
| Spanish | 80 | Translation + QA |
| French | 80 | Translation + QA |
| German | 80 | Translation + QA |
| Japanese | 120 | Translation + QA + Cultural adaptation |
| Portuguese | 80 | Translation + QA |

### Total Estimated Effort
- **Initial Setup**: 200-260 hours
- **Translation + QA**: 400-520 hours
- **Total**: 600-780 hours

### Ongoing Maintenance
- **Monthly**: 20-40 hours for updates and bug fixes
- **Quarterly**: 40-80 hours for content updates
- **Annual**: 200-300 hours for comprehensive review

### Resource Requirements
- **Linguists**: 2-3 professional translators per language
- **QA Engineers**: 1-2 for testing localization
- **Project Manager**: 1 for coordination
- **Technical Lead**: 1 for architecture and implementation

### Timeline
- **Phase 1 (Setup)**: 2-3 weeks
- **Phase 2 (Translation)**: 4-6 weeks
- **Phase 3 (Testing)**: 2-3 weeks
- **Phase 4 (Launch)**: 1-2 weeks
- **Total**: 8-14 weeks for full implementation