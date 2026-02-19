# Springfield Features & BDD

## Overview
Springfield Protocol uses Behavior Driven Development (BDD) to bridge the gap between product requirements and technical implementation. Feature files (Gherkin syntax) serve as the "Executable Truth" of the system.

## BDD Workflow

1.  **Discovery**: @Marge (Product) defines the feature requirements in `Feature.md`.
2.  **Definition**: @Lisa (Planning) translates requirements into `.feature` files in `docs/features/`.
3.  **Implementation**: @Ralph (Build) writes failing TDD tests and then implements the code to satisfy the feature specs.
4.  **Verification**: @Bart (Quality) executes the BDD scenarios to ensure acceptance criteria are met.

## Role of Agents

### @Marge (Product)
Determines "What" needs to be built. Ensures the feature solves a real user problem.

### @Bart (Quality)
The "Owner" of the feature files. Bart uses these files to perform adversarial testing, ensuring that Ralph's implementation doesn't just "work" but handles edge cases and failures gracefully.

### @Ralph (Build)
Uses the `.feature` files as the ultimate definition of "Done". Ralph's implementation is not complete until all scenarios in the feature file pass.

## Writing Features
Features are written in Gherkin syntax:

```gherkin
Feature: User Login
  As a registered user
  I want to log in
  So that I can access my dashboard

  Scenario: Successful Login
    Given I am on the login page
    When I enter valid credentials
    Then I should be redirected to the dashboard
```

## Running Tests
Springfield uses `just` to coordinate test execution.

- `just test`: Run all tests (Unit, Integration, BDD).
- `just bdd`: Run only BDD scenarios.
- `just verify`: Run tests and generate coverage report.

Example:
```bash
just bdd docs/features/login.feature
```

---
*Related: [README.md](../../README.md), [QUICK_START.md](../../QUICK_START.md)*
