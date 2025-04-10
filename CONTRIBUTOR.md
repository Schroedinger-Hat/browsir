## Continuous Integration (CI) 🛠️

Browsir uses GitHub Actions for continuous integration and automated workflows. The CI setup includes:

### Automated Checks
- **Commit Lint**: Enforces conventional commit messages (e.g., `feat:`, `fix:`, `docs:`, etc.)
- **Golangci-lint**: Runs static code analysis and style checks
- **Go Tests**: Executes all tests in the project

### Automated PR Management
- **Auto Merge**: Automatically merges PRs when they meet the following criteria:
  - PR is mergeable
  - All CI checks pass

### How to Use the CI Features

1. **Commit Messages**: Follow the conventional commit format:
   ```
   type(scope): description
   ```
   Where `type` can be:
   - `feat`: New feature
   - `fix`: Bug fix
   - `docs`: Documentation changes
   - `style`: Code style changes
   - `refactor`: Code refactoring
   - `perf`: Performance improvements
   - `test`: Adding or modifying tests
   - `build`: Build system changes
   - `ci`: CI configuration changes
   - `chore`: Maintenance tasks
   - `revert`: Reverting changes

2. **Auto Merge**: To enable automatic merging of your PR:
   - Ensure all CI checks pass
   - The PR will be automatically merged if all conditions are met

### CI Configuration Files
- `.github/workflows/ci.yaml`: Main CI workflow with linting and testing
- `.github/workflows/auto-merge.yaml`: Automated PR management workflow
- `.commitlintrc.json`: Commit message linting configuration