# Changelog

All notable changes to the goholiday project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.5] - 2025-10-04

### Changed
- **BREAKING**: Module path renamed from `github.com/coredds/GoHoliday` to `github.com/coredds/goholiday` to follow Go naming conventions
- Updated all import paths across 37 files to use lowercase naming
- Updated ChronoGo integration import path from `github.com/davidhintelmann/chronogo` to `github.com/coredds/chronogo`
- All documentation updated to reflect lowercase naming (README, API docs, CHANGELOG, CONTRIBUTING)
- Updated all example code and configuration files

### Fixed
- Corrected ChronoGo library reference to point to the correct repository
- All naming conventions now consistently use lowercase throughout the project

### Migration
Users must update their import statements:
- From: `github.com/coredds/GoHoliday` 
- To: `github.com/coredds/goholiday`

## [0.6.4] - 2025-09-29

### Added
- **Canada (CA)**: National Day for Truth and Reconciliation (September 30, since 2021)
  - Bilingual support (English/French)
  - Historical accuracy with year-based implementation
  - Proper integration with existing Canadian holidays

### Enhanced
- **Australia (AU)**: King's Birthday vs Queen's Birthday naming transition
  - Queen's Birthday for years 2022 and earlier
  - King's Birthday for years 2023 and onwards
  - Applied to both national and state-specific variations (Queensland)
  - Updated test expectations for accurate validation

- **Japan (JP)**: Tokyo Olympics holiday date adjustments (2020-2021)
  - Marine Day: July 23, 2020 and July 22, 2021 (Olympic years)
  - Mountain Day: August 10, 2020 and August 8, 2021 (Olympic years)
  - Sports Day: July 24, 2020 and July 23, 2021 (Olympic years)
  - Normal scheduling for all other years maintained

- **Netherlands (NL)**: Historical Queen's Day to King's Day transition
  - Queen's Day (April 30) for years before 2014
  - King's Day (April 27) for years 2014 and onwards
  - Proper weekend shifting rules for both holidays
  - Bilingual Dutch/English support

### Fixed
- Holiday rule accuracy aligned with Python holidays reference library
- Chronogo integration updated with all 34 supported countries
- Test coverage maintained at 100% for all updated countries
- Performance benchmarks preserved (sub-microsecond lookups)

### Technical
- Enhanced test coverage for historical holiday transitions
- Improved documentation for Olympic year adjustments
- Updated country provider tests for accuracy validation
- Maintained backward compatibility for all existing functionality

## [0.6.3] - 2025-09-18

### Added
- **Chile (CL) Holiday Provider**: Complete South American implementation with unique legal framework
  - 16 holidays including Independence Day, Army Day, Navy Day
  - Variable holidays that move to Monday per Chilean "Ley de Feriados" (Holidays Law)
  - Regional holidays: Battle of Arica (2020+), Chillán Foundation Day (2019+)
  - Easter-based holidays: Good Friday, Holy Saturday
  - Spanish/English bilingual support
  - All 16 Chilean regions supported

- **Ireland (IE) Holiday Provider**: Comprehensive Celtic and Catholic tradition implementation
  - 14 holidays including Saint Patrick's Day, bank holidays, Celtic festivals
  - Celtic seasonal festivals: Saint Brigid's Day, May Day, Lughnasadh, Samhain
  - Bank holidays on first/last Mondays of specific months
  - Saint Brigid's Day public holiday (introduced 2023)
  - Irish Gaelic/English bilingual support
  - All 30 counties and provinces supported

- **Israel (IL) Holiday Provider**: Hebrew calendar holiday system with religious and national observances
  - 11 holidays including major Jewish holidays and memorial days
  - Hebrew calendar holidays: Rosh Hashanah (2 days), Yom Kippur, Passover, Shavuot, Sukkot, Simchat Torah, Hanukkah
  - Memorial sequence: Holocaust Remembrance Day, Memorial Day, Independence Day
  - Accurate Hebrew calendar calculations for 2023-2026
  - Hebrew/English bilingual support
  - All 6 Israeli districts supported

### Enhanced
- **Multi-language Support**: Added Spanish, Irish Gaelic (GA), and Hebrew (HE) language support
- **Variable Holiday Laws**: Support for Chilean holiday movement laws and Irish bank holiday patterns
- **Hebrew Calendar Integration**: Accurate Hebrew calendar calculations for Jewish holidays
- **Regional Holiday Variations**: Enhanced support for region-specific holiday introductions and changes
- **Comprehensive Test Coverage**: 100% test coverage for all new countries with edge case validation

### Updated
- **Country Count**: Increased from 34 to 37 supported countries
- **SupportedCountries Map**: Added CL, IE, IL to validation and integration points
- **chronogo Integration**: Updated fast holiday checking for new countries
- **Configuration System**: Enhanced provider initialization for new countries
- **Sync Functionality**: Added data synchronization support for new countries
- **Documentation**: Updated README with new countries and improved examples

## [0.5.3] - 2025-08-30

### Added
- **Portugal (PT) Holiday Provider**: Complete Iberian implementation with Catholic traditions
  - 14 national holidays including Freedom Day, Portugal Day, Republic Day
  - Easter-based holidays: Carnival Tuesday, Good Friday, Easter Sunday, Corpus Christi
  - Portuguese/English bilingual support
  - 20 subdivisions (18 districts + 2 autonomous regions)
- **Italy (IT) Holiday Provider**: Comprehensive Italian holiday system with regional variations
  - 11 national holidays including Epiphany, Liberation Day, Republic Day, St. Stephen's Day
  - Regional patron saint holidays (St. Ambrose, St. Mark, etc.)
  - Italian/English bilingual support
  - All 20 Italian regions supported
- **India (IN) Holiday Provider**: Multi-religious implementation with diverse cultural celebrations
  - National holidays: Republic Day, Independence Day, Gandhi Jayanti
  - Hindu festivals: Diwali, Holi, Dussehra, Janmashtami, Ram Navami (approximate dates)
  - Buddhist, Sikh, Jain festivals: Buddha Purnima, Guru Nanak Jayanti, Mahavir Jayanti
  - State-specific holidays for major states
  - Hindi/English bilingual support
  - 36 subdivisions (28 states + 8 union territories)

### Enhanced
- **GitHub API Integration**: Added secure token authentication system
  - Token-based authentication for higher API rate limits (5000 vs 60 requests/hour)
  - Automatic token loading from environment variables or config files
  - Token validation and error handling
  - Updated sync tool with authentication support
- **Holiday Sync System**: Improved Python holidays repository integration
  - Enhanced filename mapping for new countries
  - Better error handling and rate limiting
  - Mock testing system for reliable CI/CD

### Technical
- **Test Coverage**: Added comprehensive test suites for all new countries
  - 20+ test functions covering creation, holidays, languages, categories
  - Regional/state holiday testing
  - Multi-language support validation
  - Easter calculation accuracy tests
- **Code Quality**: All new implementations follow established patterns
  - Consistent API design across all providers
  - Proper error handling and validation
  - Clean separation of concerns

## [0.5.0] - 2025-08-28

### Added
- **Norway (NO) Holiday Provider**: Complete Nordic implementation with traditional celebrations
  - Constitution Day (Grunnlovsdag), Maundy Thursday, Good Friday
  - Easter-based holidays with accurate astronomical calculations
  - Norwegian/English bilingual support
  - All 11 Norwegian counties supported
- **Turkey (TR) Holiday Provider**: Comprehensive secular and Islamic holiday implementation  
  - Democracy and National Unity Day, Victory Day, Republic Day
  - Religious holidays: Ramadan Feast, Sacrifice Feast
  - Turkish/English bilingual support
  - All 81 Turkish provinces supported
- **Russia (RU) Holiday Provider**: Orthodox calendar with extensive New Year celebrations
  - Extended New Year holidays (8-day period), Defender of the Fatherland Day
  - Orthodox Easter, Russia Day, Unity Day
  - Russian/English bilingual support
  - All 85 federal subjects supported
- **Indonesia (ID) Holiday Provider**: Multi-religious society with diverse cultural holidays
  - Islamic holidays: Eid al-Fitr, Eid al-Adha, Mawlid
  - Christian holidays: Good Friday, Ascension Day
  - Buddhist, Hindu, and Chinese New Year celebrations
  - Indonesian/English bilingual support
  - All 38 Indonesian provinces supported

### Performance
- Optimized holiday calculations for Orthodox Easter in Russia
- Enhanced multi-religious calendar support for Indonesia
- Improved subdivision lookup performance for Turkey (81 provinces)

## [0.3.0] - 2025-08-27

### Added
- **Italy (IT) Holiday Provider**: Complete implementation with comprehensive holiday coverage
  - Liberation Day, Republic Day, Assumption of Mary
  - Regional Carnival celebrations (Monday and Tuesday)
  - Easter-based holidays with accurate calculations
  - Italian/English bilingual support
  - All 20 Italian regions supported
- **Spain (ES) Holiday Provider**: Comprehensive implementation with Spanish cultural holidays
  - Epiphany, Constitution Day, Immaculate Conception
  - National Day, All Saints' Day, Good Friday
  - Spanish/English bilingual support
  - All 19 autonomous communities supported
- **Netherlands (NL) Holiday Provider**: Complete Dutch holiday implementation
  - King's Day (Koningsdag), Liberation Day, Christmas holidays
  - Easter-based holidays: Good Friday, Easter Monday, Ascension Day, Whit Monday
  - Dutch/English bilingual support
  - All 12 Dutch provinces supported
- **South Korea (KR) Holiday Provider**: Full Korean holiday implementation
  - Traditional holidays: Lunar New Year, Chuseok (Korean Thanksgiving)
  - National holidays: Independence Movement Day, Constitution Day, National Foundation Day
  - Modern holidays: Children's Day, Memorial Day, Liberation Day
  - Korean/English bilingual support
  - All 17 provinces and metropolitan cities supported

### Enhanced
- **Library Expansion**: Increased from 11 to 15 supported countries
- **Cultural Coverage**: Added 200+ regional subdivisions across new countries
- **Performance**: Maintained sub-microsecond lookup performance at 275K+ ops/sec
- **Test Coverage**: Enhanced to 82.7% in countries package with comprehensive validation
- **Repository Management**: Improved gitignore configuration for coverage artifacts

## [0.2.2] - 2025-08-27

### Added
- **Brazil (BR) Holiday Provider**: Complete implementation with 12 federal holidays
  - Carnival Monday and Tuesday with unique `carnival` category
  - Portuguese/English bilingual support
  - All 27 Brazilian states and federal district supported
  - Easter-based holidays: Good Friday, Corpus Christi
  - Cultural holidays: Independence Day, Tiradentes, All Souls' Day
- **Mexico (MX) Holiday Provider**: Comprehensive implementation with 12 holidays
  - Variable Monday holidays reflecting 2006 constitutional reforms
  - Spanish/English bilingual support  
  - All 32 Mexican federal entities supported
  - Constitutional holidays: Constitution Day, Benito Juárez Birthday, Revolution Day
  - Cultural observances: Day of the Dead, Our Lady of Guadalupe
  - Holy Week traditions: Maundy Thursday, Good Friday, Easter Saturday
- **Latin America Demo Applications**:
  - Brazil showcase with cultural highlights
  - Mexico showcase with constitutional reform explanations
  - Latin America comparison demonstrating shared and unique traditions
- **Enhanced Test Coverage**: 
  - Brazil: 9 comprehensive test functions with 100% core coverage
  - Mexico: 11 test functions covering variable holidays and cultural accuracy
  - Total of 95+ passing tests across countries package

### Improved
- **Performance Optimization**: Both new countries maintain 400K+ operations/second
  - Brazil: ~403K ops/sec (2,955 ns/op, 6,616 B/op)
  - Mexico: ~395K ops/sec (2,974 ns/op, 6,616 B/op)
- **Global Coverage**: Expanded from 9 to 11 countries (22% increase)
- **Cultural Accuracy**: Validated real-world holiday calculations and traditions
- **Code Quality**: 100% test pass rate with comprehensive linting compliance

### Technical
- **Thread Safety**: All new providers follow established concurrent patterns
- **Memory Efficiency**: Consistent allocation patterns across implementations
- **API Compatibility**: Full backward compatibility maintained
- **Integration**: Seamless integration with main goholidays.go package

### Countries Supported
- 🇺🇸 United States
- 🇨🇦 Canada
- 🇬🇧 United Kingdom
- 🇦🇺 Australia
- 🇳🇿 New Zealand
- 🇯🇵 Japan
- 🇮🇳 India
- 🇫🇷 France
- 🇩🇪 Germany
- 🇧🇷 Brazil (NEW!)
- 🇲🇽 Mexico (NEW!)

## [0.1.2] - Previous Release

### Added
- Japan (JP) Holiday Provider with 16 public holidays
- India (IN) Holiday Provider with multi-religious festivals
- France (FR) Holiday Provider with regional observances
- Enhanced chronogo integration
- Performance optimization framework
- API stability system

### Technical
- Thread-safe operations with sync.RWMutex
- Object pooling and string interning
- LRU caching implementation
- Comprehensive benchmark suite

## [0.1.1] - Earlier Release

### Added
- Core holiday providers for US, CA, GB, AU, NZ, DE
- BaseProvider pattern implementation
- Easter calculation algorithms
- Regional subdivision support

## [0.1.0] - Initial Release

### Added
- Initial goholiday library structure
- Basic holiday calculation framework
- Configuration system
- Initial test suite
