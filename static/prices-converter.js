(function () {
	const root = document.getElementById("fx-converter");
	if (!root) return;

	const rate = Number(root.dataset.rate || 0);
	if (!rate) return;

	const fromInput = document.getElementById("fx-from");
	const toInput = document.getElementById("fx-to");
	const fromLabel = document.getElementById("fx-from-label");
	const toLabel = document.getElementById("fx-to-label");
	const summary = document.getElementById("fx-summary");
	const swapBtn = document.getElementById("fx-swap");

	let direction = "usd-vnd";

	function formatVND(n) {
		return new Intl.NumberFormat("vi-VN").format(Math.round(n)) + " ₫";
	}

	function formatUSD(n) {
		return new Intl.NumberFormat("en-US", {
			minimumFractionDigits: 0,
			maximumFractionDigits: 4,
		}).format(n) + " USD";
	}

	function convert() {
		const amount = Number(fromInput.value);
		if (!Number.isFinite(amount) || amount < 0) {
			toInput.value = "";
			summary.textContent = "";
			return;
		}

		if (direction === "usd-vnd") {
			const vnd = amount * rate;
			toInput.value = Math.round(vnd);
			summary.textContent =
				formatUSD(amount) + " ≈ " + formatVND(vnd) + " (tỷ giá " + formatVND(rate) + "/USD)";
			return;
		}

		const usd = amount / rate;
		toInput.value = usd.toFixed(4).replace(/\.?0+$/, "");
		summary.textContent =
			formatVND(amount) + " ≈ " + formatUSD(usd) + " (tỷ giá " + formatVND(rate) + "/USD)";
	}

	function swap() {
		direction = direction === "usd-vnd" ? "vnd-usd" : "usd-vnd";
		const currentFrom = fromInput.value;
		const currentTo = toInput.value;

		if (direction === "usd-vnd") {
			fromLabel.textContent = "USD";
			toLabel.textContent = "VND";
			fromInput.step = "any";
			fromInput.value = currentTo || "100";
		} else {
			fromLabel.textContent = "VND";
			toLabel.textContent = "USD";
			fromInput.step = "1";
			fromInput.value = currentTo || String(rate * 100);
		}

		convert();
	}

	fromInput.addEventListener("input", convert);
	swapBtn.addEventListener("click", swap);
	convert();
})();
