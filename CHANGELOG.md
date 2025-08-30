# Changelog

All notable changes to the GoHoliday project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
- Enhanced ChronoGo integration
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
- Initial GoHoliday library structure
- Basic holiday calculation framework
- Configuration system
- Initial test suite
