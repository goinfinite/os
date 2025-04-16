# Contributing

Thank you for your interest in contributing to our project! This guide will help you understand how to effectively contribute to this repository.

## Getting Started

Before submitting any changes, please discuss your proposed modifications with the repository maintainers by opening an issue. This step is crucial to ensure that your contributions align with the project's goals and to avoid duplicating efforts.

We value meaningful contributions that substantively improve the project. Please focus on quality over quantity - contributions should address real needs or enhancement opportunities rather than superficial changes. **Proposals that appear primarily aimed at gaining contributor status without adding significant value will be declined**.

For technical assistance, please refer to the [Development](/docs/DEVELOPMENT.md) documentation. It contains information on setting up your development environment, running tests, and other essential details.

For security-related issues, please refer to our [Security Policy](/docs/SECURITY.md) for guidelines on reporting vulnerabilities and security concerns.

## Code of Conduct

Our community values respectful and inclusive collaboration. We strive to make participation a positive experience for everyone, regardless of background, identity, or experience level.

As a contributor, we ask that you:

- Treat others with respect and kindness;
- Show empathy in your interactions;
- Be receptive to constructive feedback;
- Take responsibility for mistakes and use them as learning opportunities.

Unacceptable behaviors may include, but are not limited to:

- Harassment, bullying, or intimidation;
- Discriminatory or offensive comments;
- Unwelcome attention or remarks;
- Sharing private information without permission.

If you witness or experience inappropriate behavior, please contact our community leaders at legal _at_ goinfinite.net. All reports will be reviewed and addressed appropriately.

Enforcements actions from correction to permanent bans will be taken depending on the severity of the violation.

## Licensing

This project is released under the Eclipse Public License (EPL) version 2.0. By contributing, you agree that your submissions will be governed by this license. This arrangement protects the open source nature of the project.

### Contributor Agreement

Contributors must sign the Fiduciary Contributor License Agreement (FLA) before submitting code. This legal agreement transfers copyright of contributions to a designated fiduciary who manages these rights for the project's benefit.

The FLA serves to:

- Maintain the software's free and open source status;
- Shield the project from potential copyright complications;
- Include a safety mechanism: if the fiduciary violates Free Software principles, rights return to the original contributors.

For additional information about the FLA, please refer to the [FLA FAQ](https://fsfe.org/activities/fla/fla.en.html).

## Pull Request Guidelines

When submitting pull requests, please follow these steps:

- Remove any build or installation dependencies before finalizing your build;
- Document interface changes in the README.md, including new environment variables, ports, file locations, and container parameters;
- Update version numbers according to [SemVer](http://semver.org/) conventions in relevant files (examples, CHANGELOG.md, api.go etc);
- Add unit tests for any new features or bug fixes, specially for infrastructure, value objects and any entities or use cases that contains enough non-trivial logic;
- Ensure all tests pass before submitting your pull request;
- Adhere to every Clean Code Rules described at [ntorga's article](https://ntorga.com/the-clean-coder-golden-rules/);
- Follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification for every commit message;
- Obtain approval from two developers before merging. If you lack merge permissions, the second reviewer can complete the merge for you;

By following these guidelines, you help maintain project quality and consistency while making contributions valuable to the community.
