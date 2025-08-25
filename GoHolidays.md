# Product Requirements Document (PRD) - GoHolidays

## Objective
Develop GoHolidays, a Go-native library to provide comprehensive holiday data for countries and their subdivisions, based on the well-established Python holidays library by the vacanza community. The library will bring the large body of holiday information into the Go ecosystem while following idiomatic Go design principles and ensuring high performance and usability.

## Background
The Python holidays library ([GitHub link](https://github.com/vacanza/holidays/)) efficiently generates country- and subdivision-specific government-designated holiday data in a dict-like manner. It supports 154 countries, multiple languages, different holiday categories, and can work with date strings and timestamps. This library is regularly updated and widely used in Python projects, including financial systems integration.

## Goals
- Provide accurate, up-to-date holiday information for a wide array of countries and subdivisions.
- Maintain Go idiomatic API design (avoid pythonic patterns).
- High performance and low memory footprint.
- Easily ingest updates directly from the Python holidays repository to keep country holiday data current.
- Serve as a standalone package but tightly integrate as an optional package in ChronoGo for holiday-aware date/time computations.

## Features
- **Country/Subdivision Support:** Include extensive coverage similar to Python holidays, supporting ISO country and subdivision codes.
- **Holiday Categories:** Support categories like Public, Bank, School, Government, Religious, Optional, Half Day, Armed Forces, Workday.
- **Multiple Languages:** Return holiday names in supported languages where possible.
- **Date Handling:** Accept standard date formats as input to check holidays and list holidays on specific dates.
- **Idiomatic Go API:** Use Go conventions for package structure, naming, error handling, and usage patterns.
- **Update Importing:** CLI and Go functions to easily sync and import updated holiday definitions from the Python repo (e.g. converting YAML/JSON or parsing direct Python source definitions).
- **ChronoGo Integration:** Provide extensions or hooks in ChronoGo to leverage GoHolidays for holiday-aware calculations, such as business day calculations and holiday-aware scheduling.
- **Documentation and Examples:** Well-structured documentation demonstrating installation, usage, integration with ChronoGo, and update procedures.

## Non-Goals
- Do not reimplement Python holidays logic verbatim; focus on idiomatic Go implementation.
- Do not implement a full calendar or date library, rely on ChronoGo or standard Go time package for time/date operations.
- Do not maintain holiday data manually; rely on automated update imports.

## Technical Considerations
- Use Go modules for package management.
- Package should be thread-safe.
- Design for extensibility to add new countries and categories as updates come.
- Provide a CLI tool or commands as part of the package to sync holiday data from Python repo.
- Consider efficient data structures for fast holiday lookups by date.

## Milestones
1. Initial design and API definition.
2. Core concept implementation with a subset of countries.
3. Implement update ingestion mechanism from Python holidays repo.
4. Integration with ChronoGo as optional package.
5. Documentation, testing, and release.

---

This PRD serves as a guiding document to build GoHolidays, effectively bridging the valuable holiday data from the Python ecosystem into Go for broader use cases while maintaining Go's best practices and extensibility for future updates.
