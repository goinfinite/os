"use strict";

document.addEventListener("alpine:init", () => {
  async function jsonAjax(method, url, payload) {
    const loadingOverlayElement = document.getElementById("loading-overlay");
    loadingOverlayElement.classList.add("htmx-request");

    await fetch(url, {
      method: method,
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    })
      .then((response) => {
        loadingOverlayElement.classList.remove("htmx-request");
        return response.json();
      })
      .then((parsedResponse) => {
        Alpine.store("toast").displayToast(parsedResponse.body, "success");
      })
      .catch((parsedResponse) => {
        Alpine.store("toast").displayToast(parsedResponse.body, "danger");
      });
  }

  window.jsonAjax = jsonAjax;
});
