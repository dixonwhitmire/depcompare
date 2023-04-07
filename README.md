# depcompare

depcompare compares a dependency list to a "base" dependency list. The comparison output shows dependencies
in common and dependencies which differ.

Dependency types supported include:
- [Gradle build file](internal/depcompare/build-base.gradle) aka "build.gradle"
- [Text file](internal/depcompare/gradle-base.txt) containing dependency information in the "Gradle" short format.

Usage:
```
depcompare --type=[gradleb|gradlet] [dep path] [base dep path]
```

## quickstart
```shell
make build

% ./target/depcompare --type=gradlebuild \
                      internal/depcompare/testdata/build-dep.gradle \
                      internal/depcompare/testdata/build-base.gradle
2023/04/07 16:50:45 main: dependency type set to gradlebuild
2023/04/07 16:50:45 main: loaded dependency file internal/depcompare/testdata/build-dep.gradle
2023/04/07 16:50:45 main: loaded base dependency file internal/depcompare/testdata/build-base.gradle
=======================================
Common Dependencies
=======================================
org.apache.commons:commons-collections4
org.apache.commons:commons-lang3
=======================================
=======================================
Base Only Dependencies
=======================================
org.postgresql:postgresql
org.springframework.boot:spring-boot-starter-test
org.springframework.boot:spring-boot-starter-web
=======================================
=======================================
Deps Only Dependencies
=======================================
org.apache.commons:commons-csv
=======================================
```
