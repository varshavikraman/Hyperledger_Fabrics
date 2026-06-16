const showToast = (message, type = 'success') => {
  const toast = document.getElementById("toast");
  const toastMessage = document.getElementById("toastMessage");
  const toastIcon = document.getElementById("toastIcon");

  if (!toast || !toastMessage || !toastIcon) return;

  toastMessage.textContent = message;
  toast.className = `notification show ${type}`;

  if (type === 'success') {
    toastIcon.innerHTML = `<svg width="20" height="20" fill="none" stroke="var(--success)" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`;
  } else {
    toastIcon.innerHTML = `<svg width="20" height="20" fill="none" stroke="var(--error)" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`;
  }

  setTimeout(() => {
    toast.classList.remove("show");
  }, 4500);
};

const addData = async (event) => {
  event.preventDefault();
  
  const carId = document.getElementById("carId").value.trim();
  const make = document.getElementById("make").value.trim();
  const model = document.getElementById("model").value.trim();
  const color = document.getElementById("color").value.trim();
  const dateOfManufacture = document.getElementById("dateOfManufacture").value;
  const manufacturerName = document.getElementById("manufacturerName").value.trim();

  const carData = {
    carId: carId,
    make: make,
    model: model,
    color: color,
    dateOfManufacture: dateOfManufacture,
    manufacturerName: manufacturerName,
  };

  if (!carId || !make || !model || !color || !dateOfManufacture || !manufacturerName) {
    showToast("Please enter all the details properly.", "error");
    return;
  }

  try {
    const response = await fetch("/api/car", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(carData),
    });

    console.log("RESPONSE: ", response);
    const data = await response.json();
    console.log("DATA: ", data);

    if (response.ok) {
      showToast(`Asset '${carId}' committed to ledger successfully!`, "success");
      document.getElementById("createCarForm").reset();
    } else {
      showToast(data.message || "Failed to commit asset to ledger.", "error");
    }
  } catch (err) {
    showToast("Error occurred during submission.", "error");
    console.error(err);
  }
};

const readData = async (event) => {
  event.preventDefault();
  
  const carId = document.getElementById("carIdInput").value.trim();
  const resultCard = document.getElementById("resultCard");
  const resultPlaceholder = document.getElementById("resultPlaceholder");
  const resultContent = document.getElementById("resultContent");

  if (!carId) {
    showToast("Please enter a valid ID.", "error");
    return;
  }

  try {
    const response = await fetch(`/api/car/${carId}`);
    const responseData = await response.json();
    console.log("response data:", responseData);

    const dataStr = responseData.data;
    if (!dataStr) {
      showToast("No response received from ledger", "error");
      return;
    }

    let parsedData;
    try {
      parsedData = JSON.parse(dataStr);
    } catch (e) {
      showToast("Invalid data format returned by ledger", "error");
      return;
    }

    // Handle Ledger-level errors (like asset not found)
    if (parsedData.error) {
      showToast(parsedData.error, "error");
      
      resultPlaceholder.style.display = "none";
      resultContent.style.display = "block";
      resultCard.classList.add("active");
      resultContent.innerHTML = `
        <div class="result-header" style="border-bottom-color: var(--error);">
          <span class="car-badge" style="background: var(--error);">Not Found</span>
          <span class="result-car-title">Query Error</span>
        </div>
        <div style="padding: 1rem 0; color: var(--text-secondary); font-size: 0.875rem;">
          <p>${parsedData.error}</p>
        </div>
      `;
      return;
    }

    // Render successful query details
    showToast("Asset details fetched successfully!", "success");
    
    resultPlaceholder.style.display = "none";
    resultContent.style.display = "block";
    resultCard.classList.add("active");
    
    resultContent.innerHTML = `
      <div class="result-header">
        <span class="car-badge">Verified</span>
        <span class="result-car-title">${parsedData.make} ${parsedData.model}</span>
      </div>
      <ul class="spec-list">
        <li class="spec-item">
          <span class="spec-label">Asset ID</span>
          <span class="spec-value" style="color: var(--accent-cyan); font-family: monospace;">${parsedData.carId}</span>
        </li>
        <li class="spec-item">
          <span class="spec-label">Color</span>
          <span class="spec-value">${parsedData.color}</span>
        </li>
        <li class="spec-item">
          <span class="spec-label">Manufacture Date</span>
          <span class="spec-value">${parsedData.dateOfManufacture}</span>
        </li>
        <li class="spec-item">
          <span class="spec-label">Manufacturer Name</span>
          <span class="spec-value">${parsedData.manufacturerName}</span>
        </li>
      </ul>
    `;

  } catch (err) {
    showToast("Error fetching from ledger.", "error");
    console.error(err);
  }
};

// Bind to window scope to allow execution from HTML onsubmit events
window.addData = addData;
window.readData = readData;