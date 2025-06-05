UiToolset.RegisterAlpineState(() => {
  Alpine.data("runtimes", () => ({
    // PrimaryState
    selectedVhostHostname: "",
    selectedRuntimeType: "",
    vhostHostname: selectedVhostHostname,
    reloadRuntimePageContent(vhostHostname, runtimeType) {
      htmx.ajax(
        "GET",
        "/runtimes/?vhostHostname=" +
          vhostHostname +
          "&runtimeType=" +
          runtimeType,
        {
          select: "#runtimes-page-content",
          target: "#runtimes-page-content",
          indicator: "#loading-overlay",
          swap: "outerHTML transition:true",
        }
      );
    },
    init() {
      this.selectedVhostHostname = document.getElementById(
        "selectedVhostHostname"
      ).value;
      this.selectedRuntimeType = document.getElementById(
        "selectedRuntimeType"
      ).value;
    },

    updateSelectedVhostHostname(vhostHostname) {
      this.reloadRuntimePageContent(vhostHostname, this.selectedRuntimeType);
    },
    updateSelectedRuntimeType(runtimeType) {
      this.reloadRuntimePageContent(this.selectedVhostHostname, runtimeType);
    },
  }));
});
