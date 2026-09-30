const loginForm = document.querySelector("#login-form");
const loginMessage = document.querySelector("#login-message");
const logoutButton = document.querySelector("#logout-button");

const form = document.querySelector("#request-form");
const message = document.querySelector("#message");

const requestsBox = document.querySelector("#requests");
const refreshButton = document.querySelector("#refresh");

let token = localStorage.getItem("domservice_token") || "";

loginForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  loginMessage.textContent = "Logging in...";

  const payload = {
    email: document.querySelector("#login-email").value,
    password: document.querySelector("#login-password").value,
  };

  try {
    const response = await fetch("/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    const data = await response.json();

    if (!response.ok) {
      loginMessage.textContent = data.error || "Login failed";
      return;
    }

    token = data.token;

    localStorage.setItem("domservice_token", token);

    loginMessage.textContent =
      `Logged in as ${data.user.email} (${data.user.role})`;

    await loadRequests();
  } catch (error) {
    loginMessage.textContent = "Could not connect to server";
  }
});

logoutButton.addEventListener("click", () => {
  token = "";

  localStorage.removeItem("domservice_token");

  loginMessage.textContent = "Logged out";

  requestsBox.innerHTML =
    "<p class='empty'>Login as admin to see requests.</p>";
});

form.addEventListener("submit", async (event) => {
  event.preventDefault();

  message.textContent = "Saving...";

  const payload = {
    name: document.querySelector("#name").value,
    phone: document.querySelector("#phone").value,
    address: document.querySelector("#address").value,
    problem: document.querySelector("#problem").value,
  };

  try {
    const response = await fetch("/api/requests", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    const data = await response.json();

    if (!response.ok) {
      message.textContent = data.error || "Something went wrong";
      return;
    }

    form.reset();

    message.textContent =
      `Request #${data.id} created successfully`;

    if (token) {
      await loadRequests();
    }
  } catch (error) {
    message.textContent = "Could not connect to server";
  }
});

refreshButton.addEventListener("click", loadRequests);

async function loadRequests() {
  if (!token) {
    requestsBox.innerHTML =
      "<p class='empty'>Login as admin to see requests.</p>";
    return;
  }

  requestsBox.innerHTML =
    "<p class='empty'>Loading...</p>";

  try {
    const response = await fetch("/api/requests", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (response.status === 401 || response.status === 403) {
      requestsBox.innerHTML =
        "<p class='empty'>Admin authorization required.</p>";
      return;
    }

    if (!response.ok) {
      requestsBox.innerHTML =
        "<p class='empty'>Could not load requests.</p>";
      return;
    }

    const items = await response.json();

    if (!items || items.length === 0) {
      requestsBox.innerHTML =
        "<p class='empty'>No requests yet.</p>";
      return;
    }

    requestsBox.innerHTML = items.map((item) => `
      <article class="request-item">

        <div class="request-title">
          <strong>
            #${item.id} · ${escapeHTML(item.name)}
          </strong>

          <span class="badge">
            ${escapeHTML(item.status)}
          </span>
        </div>

        <p>
          ${escapeHTML(item.problem)}
        </p>

        <p class="meta">
          ${escapeHTML(item.phone)}
          ·
          ${escapeHTML(item.address)}
        </p>

        <div class="request-actions">

          <select onchange="changeStatus(${item.id}, this.value)">
            ${statusOption("new", item.status)}
            ${statusOption("in_progress", item.status)}
            ${statusOption("done", item.status)}
            ${statusOption("cancelled", item.status)}
          </select>

          <button
            class="danger"
            onclick="deleteRequest(${item.id})"
          >
            Delete
          </button>

        </div>

      </article>
    `).join("");
  } catch (error) {
    requestsBox.innerHTML =
      "<p class='empty'>Could not connect to server.</p>";
  }
}

function statusOption(value, current) {
  return `
    <option
      value="${value}"
      ${value === current ? "selected" : ""}
    >
      ${value}
    </option>
  `;
}

async function changeStatus(id, status) {
  if (!token) {
    alert("Please login first");
    return;
  }

  const response = await fetch(
    `/api/requests/${id}/status`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        status: status,
      }),
    }
  );

  if (!response.ok) {
    alert("Could not update status");
    return;
  }

  await loadRequests();
}

async function deleteRequest(id) {
  if (!token) {
    alert("Please login first");
    return;
  }

  const confirmed = confirm(
    `Delete request #${id}?`
  );

  if (!confirmed) {
    return;
  }

  const response = await fetch(
    `/api/requests/${id}`,
    {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  if (!response.ok) {
    alert("Could not delete request");
    return;
  }

  await loadRequests();
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

if (token) {
  loginMessage.textContent =
    "Saved login token found";

  loadRequests();
} else {
  requestsBox.innerHTML =
    "<p class='empty'>Login as admin to see requests.</p>";
}