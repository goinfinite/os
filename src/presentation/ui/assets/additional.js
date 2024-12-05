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

  function createRandomPassword() {
		const passwordLength = 16;
		const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+';

		let passwordContent = '';
		let passwordIterationCount = 0;
		while (passwordIterationCount < passwordLength) {
			const randomIndex = Math.floor(Math.random() * chars.length);
			const indexAfterRandomIndex = randomIndex + 1;
			passwordContent += chars.substring(randomIndex, indexAfterRandomIndex);

			passwordIterationCount++;
		}

		return passwordContent;
	}

	window.Infinite = {
		JsonAjax: jsonAjax,
		CreateRandomPassword: createRandomPassword
	}
});
