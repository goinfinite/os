"use strict";

document.addEventListener("alpine:init", () => {
  async function jsonAjax(method, url, payload) {
    const loadingOverlayElement = document.getElementById("loading-overlay");
    loadingOverlayElement.classList.add("htmx-request");

	try {
		const response = await fetch(url, {
			method: method,
			headers: {
				Accept: "application/json",
				"Content-Type": "application/json",
			},
			body: JSON.stringify(payload),
		});
		const parsedResponse = await response.json()

		loadingOverlayElement.classList.remove("htmx-request");

		if (!response.ok) {
			throw new Error(parsedResponse.body);
		}

		if (method.toUpperCase() !== "GET") {
			Alpine.store("toast").displayToast(parsedResponse.body, "success");
		}

		return parsedResponse.body;
	} catch (error) {
		Alpine.store("toast").displayToast(error.message, "danger");
	}
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
