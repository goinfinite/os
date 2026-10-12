UiToolset.RegisterAlpineState(() => {
  Alpine.data("collapsibleText", () => ({
    isTextContentBiggerThanAllowed: false,
    shouldDisplayPartialContent: true,
    originalTextContent: "",
    maxCharsToDisplay: 0,

    get truncatedTextContent() {
      if (
        this.shouldDisplayPartialContent &&
        this.isTextContentBiggerThanAllowed
      ) {
        return `${this.originalTextContent.substring(0, this.maxCharsToDisplay)}...`;
      }
      return this.originalTextContent;
    },

    init() {
      this.originalTextContent = this.$el.dataset.textContent ?? "";
      this.maxCharsToDisplay = Number.parseInt(
        this.$el.dataset.maxCharsToDisplay ?? "0",
        10,
      );
      this.isTextContentBiggerThanAllowed =
        this.originalTextContent.length > this.maxCharsToDisplay;
    },

    toggleTextContentDisplay() {
      this.shouldDisplayPartialContent = !this.shouldDisplayPartialContent;
    },
  }));
});
