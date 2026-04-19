const API_BASE = "/api";
const USER_TOKEN_KEY = "parts_user_token";

function showToast(message) {
  const toast = document.getElementById("toast");
  if (!toast) {
    return;
  }
  toast.textContent = message;
  toast.classList.add("toast--visible");
  setTimeout(() => toast.classList.remove("toast--visible"), 2300);
}

function setFormError(id, message) {
  const node = document.getElementById(id);
  if (!node) {
    showToast(message);
    return;
  }
  node.textContent = message || "";
  node.classList.toggle("is-visible", Boolean(message));
}

function clearFieldErrors(form) {
  if (!form) return;
  form.querySelectorAll(".field-error").forEach((n) => n.classList.remove("field-error"));
}

function markFieldError(input) {
  if (!input) return;
  input.classList.add("field-error");
  input.focus?.();
}

function getUserToken() {
  return localStorage.getItem(USER_TOKEN_KEY) || "";
}

function setUserToken(token) {
  localStorage.setItem(USER_TOKEN_KEY, token);
  updateSessionLabel();
}

function updateSessionLabel() {
  const node = document.getElementById("session-label");
  if (!node) return;
  if (getUserToken()) {
    node.textContent = "Сессия: покупатель";
    return;
  }
  node.textContent = "Гость";
}

function money(value) {
  return `${new Intl.NumberFormat("ru-RU").format(value)} ₽`;
}

function isValidEmail(value) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(value || "").trim());
}

function normalizePhone(value) {
  return String(value || "").replace(/[^\d+]/g, "").trim();
}

async function readErrorMessage(response) {
  const ct = (response.headers.get("content-type") || "").toLowerCase();
  try {
    if (ct.includes("application/json")) {
      const data = await response.json();
      if (data && typeof data === "object") {
        if (typeof data.error === "string" && data.error.trim()) {
          return data.error.trim();
        }
        if (typeof data.message === "string" && data.message.trim()) {
          return data.message.trim();
        }
      }
    } else {
      const text = (await response.text()).trim();
      if (text) {
        return text.length > 280 ? `${text.slice(0, 280)}…` : text;
      }
    }
  } catch (_) {
    // ignore
  }

  const fallback = response.statusText || "Request failed";
  return `${response.status} ${fallback}`.trim();
}

/** Бэкенд иногда отдаёт массив, иногда строку-сообщение, иногда { items: [...] } */
function normalizeListResponse(data) {
  if (Array.isArray(data)) {
    return { items: data, hint: "" };
  }
  if (typeof data === "string") {
    return { items: [], hint: data };
  }
  if (data && typeof data === "object" && Array.isArray(data.items)) {
    return { items: data.items, hint: typeof data.message === "string" ? data.message : "" };
  }
  return { items: [], hint: "" };
}

async function apiRequest(path, options = {}, role = "user") {
  const headers = { ...(options.headers || {}) };
  const token = getUserToken();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return fetch(`${API_BASE}${path}`, { ...options, headers });
}

function renderParts(items) {
  const node = document.getElementById("parts-result");
  if (!items.length) {
    node.innerHTML = "<p class='muted'>По вашему запросу нет результатов.</p>";
    return;
  }
  node.innerHTML = items.map((item) => {
    const encoded = encodeURIComponent(JSON.stringify(item));
    return `
    <article class="item">
      <div class="item__row">
        <div>
          <strong>${item.name}</strong>
          <p class="muted">${item.brand} • ${item.part_id} • ${item.delivery_day} дн.</p>
          <p>${money(item.price)}</p>
        </div>
        <button type="button" class="btn btn--primary" data-add-part="${encoded}">Добавить</button>
      </div>
    </article>
  `;
  }).join("");
}

async function loadCart() {
  const node = document.getElementById("cart-container");
  const response = await apiRequest("/user/cart");
  if (!response.ok) {
    node.innerHTML = "<p class='muted'>Войдите как покупатель для работы с корзиной.</p>";
    return;
  }
  const data = await response.json();
  const { items, hint } = normalizeListResponse(data);
  if (!items.length) {
    node.innerHTML = `<p class='muted'>${hint || "Корзина пуста."}</p>`;
    return;
  }
  node.innerHTML = `
    ${items.map((item) => `
      <article class="item">
        <div class="item__row">
          <div>
            <strong>${item.name}</strong>
            <p class="muted">${item.brand} • ${money(item.price)} × ${item.quantity}</p>
          </div>
          <button class="btn btn--danger" data-remove-part="${encodeURIComponent(item.part_id)}">Удалить</button>
        </div>
      </article>
    `).join("")}
    <p><strong>Итого: ${money(data.total || 0)}</strong></p>
  `;
}

async function loadCars() {
  const node = document.getElementById("cars-list");
  const response = await apiRequest("/user/cars");
  if (!response.ok) {
    node.innerHTML = "<p class='muted'>Авторизуйтесь как покупатель.</p>";
    return;
  }
  const data = await response.json();
  const { items, hint } = normalizeListResponse(data);
  if (!items.length) {
    node.innerHTML = `<p class='muted'>${hint || "Нет добавленных авто."}</p>`;
    return;
  }
  node.innerHTML = items.map((item) => `
    <article class="item">
      <div class="item__row">
        <div>
          <strong>${item.name}</strong>
          <p class="muted">VIN: ${item.vin}</p>
        </div>
        <button class="btn btn--danger" data-car-delete="${item.id}">Удалить</button>
      </div>
    </article>
  `).join("");
}

function renderUserOrders(items) {
  const node = document.getElementById("user-orders");
  if (!items.length) {
    node.innerHTML = "<p class='muted'>У вас пока нет заказов.</p>";
    return;
  }
  node.innerHTML = items.map((order) => `
    <article class="item">
      <strong>Заказ #${order.id}</strong>
      <p class="muted">Статус: ${order.status} • Оплата: ${order.payment_status}</p>
      <p class="muted">Адрес: ${order.address}</p>
      <p><strong>${money(order.total || 0)}</strong></p>
      ${order.payment_url ? `<a class="btn" href="${order.payment_url}" target="_blank" rel="noopener">Ссылка на оплату</a>` : ""}
    </article>
  `).join("");
}

async function loadUserOrders() {
  const response = await apiRequest("/user/orders");
  if (!response.ok) {
    renderUserOrders([]);
    return;
  }
  const data = await response.json();
  const { items } = normalizeListResponse(data);
  renderUserOrders(items);
}

document.getElementById("search-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const query = document.getElementById("search-input").value.trim();
  if (!query) {
    return;
  }
  const response = await fetch(`${API_BASE}/parts/search?q=${encodeURIComponent(query)}`);
  if (!response.ok) {
    showToast("Ошибка поиска");
    return;
  }
  const data = await response.json();
  renderParts(data.items || []);
});

document.getElementById("parts-result").addEventListener("click", async (event) => {
  const button = event.target.closest("[data-add-part]");
  if (!button) {
    return;
  }
  if (!getUserToken()) {
    showToast("Сначала войдите как покупатель");
    return;
  }
  const payload = JSON.parse(decodeURIComponent(button.getAttribute("data-add-part") || ""));
  payload.quantity = 1;
  payload.image_url = "";
  const response = await apiRequest("/user/cart/items", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });
  if (!response.ok) {
    showToast("Не удалось добавить в корзину");
    return;
  }
  showToast("Товар добавлен в корзину");
  loadCart();
});

document.getElementById("cart-container").addEventListener("click", async (event) => {
  const button = event.target.closest("[data-remove-part]");
  if (!button) {
    return;
  }
  const response = await apiRequest(`/user/cart/items/${button.dataset.removePart}`, { method: "DELETE" });
  if (!response.ok) {
    showToast("Не удалось удалить товар");
    return;
  }
  showToast("Позиция удалена");
  loadCart();
});

document.getElementById("checkout-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const statusNode = document.getElementById("checkout-status");
  const payload = { address: document.getElementById("checkout-address").value.trim() };
  if (!payload.address) {
    statusNode.textContent = "Введите адрес доставки.";
    return;
  }
  const response = await apiRequest("/user/cart/checkout", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });
  if (!response.ok) {
    statusNode.textContent = await readErrorMessage(response);
    return;
  }
  const data = await response.json();
  const orderId = data.order_id ?? data.id;
  const amount = data.amount ?? data.total ?? 0;
  statusNode.textContent = orderId
    ? `Заказ #${orderId} создан.${amount ? ` Сумма: ${money(amount)}.` : ""} ${data.message || ""}`.trim()
    : (data.message || "Заказ создан.");
  showToast("Заказ оформлен");
  if (data.payment_url) {
    window.open(data.payment_url, "_blank", "noopener");
  }
  loadCart();
  loadUserOrders();
});

document.getElementById("register-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  clearFieldErrors(form);
  setFormError("register-error", "");

  const submitBtn = form.querySelector('button[type="submit"]');
  if (submitBtn) submitBtn.disabled = true;

  try {
    const payload = {
      name: document.getElementById("reg-name").value.trim(),
      phone: normalizePhone(document.getElementById("reg-phone").value),
      email: document.getElementById("reg-email").value.trim(),
      password: document.getElementById("reg-password").value,
      address: document.getElementById("reg-address").value.trim()
    };

    if (!payload.name) {
      setFormError("register-error", "Введите имя");
      markFieldError(document.getElementById("reg-name"));
      return;
    }
    if (!payload.phone) {
      setFormError("register-error", "Введите телефон");
      markFieldError(document.getElementById("reg-phone"));
      return;
    }
    if (!isValidEmail(payload.email)) {
      setFormError("register-error", "Введите корректный email");
      markFieldError(document.getElementById("reg-email"));
      return;
    }
    if (String(payload.password || "").trim().length < 6) {
      setFormError("register-error", "Пароль должен быть минимум 6 символов");
      markFieldError(document.getElementById("reg-password"));
      return;
    }
    if (!payload.address) {
      setFormError("register-error", "Введите адрес");
      markFieldError(document.getElementById("reg-address"));
      return;
    }

    const response = await fetch(`${API_BASE}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    if (!response.ok) {
      setFormError("register-error", await readErrorMessage(response));
      return;
    }

    // Register returns a user model; obtain token via login.
    const loginResponse = await fetch(`${API_BASE}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: payload.email, password: payload.password })
    });
    if (!loginResponse.ok) {
      setFormError("register-error", await readErrorMessage(loginResponse));
      return;
    }
    const tokens = await loginResponse.json();
    setUserToken(tokens.access_token);
    showToast("Покупатель зарегистрирован и авторизован");

    loadCars();
    loadCart();
    loadUserOrders();
  } catch (err) {
    const msg = err && err.message ? err.message : String(err);
    setFormError("register-error", `Ошибка: ${msg}`);
    showToast(`Ошибка: ${msg}`);
  } finally {
    if (submitBtn) submitBtn.disabled = false;
  }
});

document.getElementById("login-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  clearFieldErrors(form);
  setFormError("login-error", "");

  const submitBtn = form.querySelector('button[type="submit"]');
  if (submitBtn) submitBtn.disabled = true;

  try {
    const payload = {
      email: document.getElementById("login-email").value.trim(),
      password: document.getElementById("login-password").value
    };

    if (!isValidEmail(payload.email)) {
      setFormError("login-error", "Введите корректный email");
      markFieldError(document.getElementById("login-email"));
      return;
    }
    if (!payload.password) {
      setFormError("login-error", "Введите пароль");
      markFieldError(document.getElementById("login-password"));
      return;
    }

    const response = await fetch(`${API_BASE}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });

    if (!response.ok) {
      setFormError("login-error", await readErrorMessage(response));
      return;
    }

    const data = await response.json();

    setUserToken(data.access_token); // ✅ FIX
    showToast("Покупатель авторизован");

    loadCars();
    loadCart();
    loadUserOrders();
  } catch (err) {
    const msg = err && err.message ? err.message : String(err);
    setFormError("login-error", `Ошибка: ${msg}`);
    showToast(`Ошибка: ${msg}`);
  } finally {
    if (submitBtn) submitBtn.disabled = false;
  }
});

document.getElementById("car-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const payload = {
    name: document.getElementById("car-name").value.trim(),
    vin: document.getElementById("car-vin").value.trim()
  };
  const response = await apiRequest("/user/cars", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });
  if (!response.ok) {
    showToast("Ошибка при добавлении авто");
    return;
  }
  showToast("Авто добавлено");
  loadCars();
});

document.getElementById("cars-list").addEventListener("click", async (event) => {
  const button = event.target.closest("[data-car-delete]");
  if (!button) {
    return;
  }
  const response = await apiRequest(`/user/cars/${button.dataset.carDelete}`, { method: "DELETE" });
  if (!response.ok) {
    showToast("Ошибка удаления авто");
    return;
  }
  showToast("Авто удалено");
  loadCars();
});

updateSessionLabel();
loadCars();
loadCart();
loadUserOrders();
