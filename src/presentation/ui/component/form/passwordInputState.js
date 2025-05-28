Infinite.RegisterAlpineState(passwordInputAlpineState);

function passwordInputAlpineState() {
  Alpine.data("passwordInput", () => ({
    // RandomPasswordGeneratorState
    isPasswordReadable: false,
    generateRandomPassword() {
      const passwordContent = Infinite.CreateRandomPassword();

      this.displayPasswordStrengthCriteria = false;
      this.updatePasswordStrengthPercentage(passwordContent);

      return passwordContent;
    },
    // Password Strength Criteria States
    displayPasswordStrengthCriteria: false,
    passwordStrengthPercentage: 0,
    passwordStrengthCriteria: {},
    resetPasswordStrengthParams() {
      this.passwordStrengthPercentage = 0;
      this.passwordStrengthCriteria = {
        isLongEnough: false,
        hasNumbers: false,
        hasUppercaseChars: false,
        hasLowercaseChars: false,
        hasSpecialChars: false,
      };
    },
    updatePasswordStrengthPercentage(password) {
      this.resetPasswordStrengthParams();

      let passwordStrengthPercentage = 0;
      if (password.length >= 6 && password.length <= 64) {
        this.passwordStrengthCriteria.isLongEnough = true;
        passwordStrengthPercentage += 20;
      }

      if (/[1-9]/.test(password)) {
        this.passwordStrengthCriteria.hasNumbers = true;
        passwordStrengthPercentage += 20;
      }

      if (/[A-Z]/.test(password)) {
        this.passwordStrengthCriteria.hasUppercaseChars = true;
        passwordStrengthPercentage += 20;
      }

      if (/[a-z]/.test(password)) {
        this.passwordStrengthCriteria.hasLowercaseChars = true;
        passwordStrengthPercentage += 20;
      }

      if (/[!@#\$%\^\&*\)\(+=._-]/.test(password)) {
        this.passwordStrengthCriteria.hasSpecialChars = true;
        passwordStrengthPercentage += 20;
      }

      this.passwordStrengthPercentage = passwordStrengthPercentage;
    },
  }));
}
