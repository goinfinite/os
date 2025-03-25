document.addEventListener("alpine:init", () => {
  const selectedVhostHostname = document.getElementById(
    "selectedVhostHostname"
  ).value;
  const selectedRuntimeType = document.getElementById(
    "selectedRuntimeType"
  ).value;

  Alpine.data("runtimes", () => ({
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
    updateSelectedVhostHostname(vhostHostname) {
      this.reloadRuntimePageContent(vhostHostname, selectedRuntimeType);
    },
    updateSelectedRuntimeType(runtimeType) {
      this.reloadRuntimePageContent(selectedVhostHostname, runtimeType);
    },
  }));
});
