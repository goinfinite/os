UiToolset.RegisterAlpineState(() => {
  Alpine.data("runtimes", () => ({
    // PrimaryState
    vhostHostname: "",
    runtimeType: "",
    reloadRuntimePageContent(vhostHostname, runtimeType) {
      htmx.ajax(
        "GET",
        document.baseURI +
          "runtimes/?vhostHostname=" +
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
      this.vhostHostname = document.getElementById("vhostHostname").value;
      this.runtimeType = document.getElementById("runtimeType").value;
    },

    updateVhostHostname(vhostHostname) {
      this.reloadRuntimePageContent(vhostHostname, this.runtimeType);
    },
    updateRuntimeType(runtimeType) {
      this.reloadRuntimePageContent(this.vhostHostname, runtimeType);
    },
  }));
});
