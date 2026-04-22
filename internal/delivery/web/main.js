const API_BASE = "/api";
const USER_TOKEN_KEY = "parts_user_token";
const USER_PROFILE_KEY = "parts_user_profile";

const APP_PAGES = ["search", "cabinet", "cart"];

function showToast(message) {
  const toast = document.getElementById("toast");
  if (!toast) return;
  toast.textContent = message;
  toast.classList.add("toast--visible");
  setTimeout(() => toast.classList.remove("toast--visible"), 2300);
}

function escapeHtml(s) {
  return String(s ?? "")
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
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

function getProfile() {
  try {
    return JSON.parse(localStorage.getItem(USER_PROFILE_KEY) || "{}");
  } catch {
    return {};
  }
}

function saveProfile(partial) {
  const prev = getProfile();
  localStorage.setItem(USER_PROFILE_KEY, JSON.stringify({ ...prev, ...partial }));
}

function parseAccessTokenUserId(token) {
  if (!token) return null;
  try {
    const parts = token.split(".");
    if (parts.length < 2) return null;
    const b64 = parts[1].replace(/-/g, "+").replace(/_/g, "/");
    const pad = b64.length % 4;
    const padded = pad ? b64 + "=".repeat(4 - pad) : b64;
    const json = atob(padded);
    const payload = JSON.parse(json);
    if (typeof payload.user_id === "number") return payload.user_id;
    if (typeof payload.user_id === "string") return parseInt(payload.user_id, 10) || null;
    return null;
  } catch {
    return null;
  }
}

function setUserToken(token) {
  localStorage.setItem(USER_TOKEN_KEY, token);
  updateSessionLabel();
  applyRoute();
}

function clearSession() {
  localStorage.removeItem(USER_TOKEN_KEY);
  localStorage.removeItem(USER_PROFILE_KEY);
  updateSessionLabel();
  applyRoute();
}

function updateSessionLabel() {
  const node = document.getElementById("session-label");
  if (!node) return;
  if (getUserToken()) {
    node.textContent = "В системе";
    return;
  }
  node.textContent = "Гость";
}

function renderProfilePanel() {
  const p = getProfile();
  const uid = parseAccessTokenUserId(getUserToken());
  const elName = document.getElementById("profile-name");
  const elEmail = document.getElementById("profile-email");
  const elPhone = document.getElementById("profile-phone");
  const elUid = document.getElementById("profile-user-id");
  if (elName) elName.textContent = p.name || "—";
  if (elEmail) elEmail.textContent = p.email || "—";
  if (elPhone) elPhone.textContent = p.phone || "—";
  if (elUid) elUid.textContent = uid != null ? String(uid) : "—";
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
    /* ignore */
  }
  const fallback = response.statusText || "Request failed";
  return `${response.status} ${fallback}`.trim();
}

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

async function apiRequest(path, options = {}) {
  const headers = { ...(options.headers || {}) };
  const token = getUserToken();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return fetch(`${API_BASE}${path}`, { ...options, headers });
}

function parseHashPage() {
  const raw = (location.hash || "#/").replace(/^#\/?/, "").split("/")[0] || "search";
  return APP_PAGES.includes(raw) ? raw : "search";
}

function setActiveAppPage(page) {
  document.querySelectorAll(".app-page").forEach((el) => {
    el.classList.toggle("is-active", el.dataset.page === page);
  });
  document.querySelectorAll(".app-nav__link").forEach((el) => {
    el.classList.toggle("is-active", el.dataset.page === page);
  });
  if (page === "cabinet") {
    renderProfilePanel();
    loadCars();
    loadUserOrders();
  }
  if (page === "cart") {
    loadCart();
  }
}

function applyRoute() {
  const viewAuth = document.getElementById("view-auth");
  const viewApp = document.getElementById("view-app");
  const appNav = document.getElementById("app-nav");
  const token = getUserToken();

  if (!token) {
    viewAuth?.classList.remove("is-hidden");
    viewApp?.classList.add("is-hidden");
    appNav?.classList.add("is-hidden");
    if (location.hash && location.hash !== "#/" && location.hash !== "") {
      history.replaceState(null, "", `${location.pathname}${location.search}#/`);
    }
    return;
  }

  viewAuth?.classList.add("is-hidden");
  viewApp?.classList.remove("is-hidden");
  appNav?.classList.remove("is-hidden");

  let page = parseHashPage();
  if (!location.hash || location.hash === "#/" || location.hash === "#") {
    page = "search";
    history.replaceState(null, "", `${location.pathname}${location.search}#/search`);
  }
  setActiveAppPage(page);
}

function switchAuthTab(tab) {
  const loginTab = document.getElementById("tab-login");
  const regTab = document.getElementById("tab-register");
  const loginPanel = document.getElementById("panel-login");
  const regPanel = document.getElementById("panel-register");
  const isLogin = tab === "login";
  loginTab?.classList.toggle("is-active", isLogin);
  regTab?.classList.toggle("is-active", !isLogin);
  loginTab?.setAttribute("aria-selected", isLogin ? "true" : "false");
  regTab?.setAttribute("aria-selected", isLogin ? "false" : "true");
  loginPanel?.classList.toggle("is-active", isLogin);
  regPanel?.classList.toggle("is-active", !isLogin);
  setFormError("login-error", "");
  setFormError("register-error", "");
}

function renderParts(items) {
  const node = document.getElementById("parts-result");
  if (!node) return;
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
          <strong>${escapeHtml(item.name)}</strong>
          <p class="muted">${escapeHtml(item.brand)} • ${escapeHtml(item.part_id)} • ${escapeHtml(String(item.delivery_day))} дн.</p>
          <p>${money(item.price)}</p>
        </div>
        <div>
          <button type="button" class="btn btn--primary" data-add-part="${encoded}">В корзину</button>
          <button type="button" class="btn" data-check-part="${escapeHtml(item.part_id)}" data-check-name="${escapeHtml(item.name)}">Уточнить срок</button>
        </div>
      </div>
    </article>
  `;
  }).join("");
}

async function loadCart() {
  const node = document.getElementById("cart-container");
  if (!node) return;
  const response = await apiRequest("/user/cart");
  if (!response.ok) {
    node.innerHTML = "<p class='muted'>Не удалось загрузить корзину.</p>";
    return;
  }
  const data = await response.json();
  const { items, hint } = normalizeListResponse(data);
  if (!items.length) {
    node.innerHTML = `<p class='muted'>${escapeHtml(hint || "Корзина пуста.")}</p>`;
    return;
  }
  node.innerHTML = `
    ${items.map((item) => `
      <article class="item">
        <div class="item__row">
          <div>
            <strong>${escapeHtml(item.name)}</strong>
            <p class="muted">${escapeHtml(item.brand)} • ${money(item.price)} × ${escapeHtml(String(item.quantity))}</p>
          </div>
          <button type="button" class="btn btn--danger" data-remove-part="${encodeURIComponent(item.part_id)}">Удалить</button>
        </div>
      </article>
    `).join("")}
    <p><strong>Итого: ${money(data.total || 0)}</strong></p>
  `;
}

function carVin(item) {
  if (item == null) return "";
  return item.vin ?? item.VIN ?? item.Vin ?? "";
}

async function loadCars() {
  const node = document.getElementById("cars-list");
  if (!node) return;
  const response = await apiRequest("/user/cars");
  if (!response.ok) {
    node.innerHTML = "<p class='muted'>Не удалось загрузить автомобили.</p>";
    return;
  }
  const data = await response.json();
  const { items, hint } = normalizeListResponse(data);
  if (!items.length) {
    node.innerHTML = `<p class='muted'>${escapeHtml(hint || "Нет добавленных авто.")}</p>`;
    return;
  }
  node.innerHTML = items.map((item) => {
    const vin = carVin(item);
    return `
    <article class="item">
      <div class="item__row">
        <div>
          <strong>${escapeHtml(item.name)}</strong>
          <p class="muted vin-line">VIN: <span class="vin-code">${escapeHtml(vin)}</span></p>
        </div>
        <button type="button" class="btn btn--danger" data-car-delete="${escapeHtml(String(item.id))}">Удалить</button>
      </div>
    </article>
  `;
  }).join("");
}

function renderUserOrders(items) {
  const node = document.getElementById("user-orders");
  if (!node) return;
  if (!items.length) {
    node.innerHTML = "<p class='muted'>У вас пока нет заказов.</p>";
    return;
  }
  node.innerHTML = items.map((order) => `
    <article class="item">
      <strong>Заказ #${escapeHtml(String(order.id))}</strong>
      <p class="muted">Статус: ${escapeHtml(String(order.status))} • Оплата: ${escapeHtml(String(order.payment_status))}</p>
      <p class="muted">Адрес: ${escapeHtml(String(order.address))}</p>
      <p><strong>${money(order.total || 0)}</strong></p>
      ${order.payment_url ? `<a class="btn" href="${escapeHtml(order.payment_url)}" target="_blank" rel="noopener">Ссылка на оплату</a>` : ""}
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

function afterAuthSuccess() {
  renderProfilePanel();
  loadCars();
  loadCart();
  loadUserOrders();
  if (!location.hash || location.hash === "#/" || location.hash === "#") {
    history.replaceState(null, "", `${location.pathname}${location.search}#/search`);
  }
  applyRoute();
}

document.getElementById("tab-login")?.addEventListener("click", () => switchAuthTab("login"));
document.getElementById("tab-register")?.addEventListener("click", () => switchAuthTab("register"));

window.addEventListener("hashchange", () => applyRoute());

document.getElementById("btn-logout")?.addEventListener("click", () => {
  clearSession();
  switchAuthTab("login");
  showToast("Вы вышли из аккаунта");
});

document.getElementById("search-form")?.addEventListener("submit", async (event) => {
  event.preventDefault();
  const query = document.getElementById("search-input")?.value?.trim();
  const statusNode = document.getElementById("search-status");
  if (!query) {
    showToast("Введите запрос");
    return;
  }
  const response = await apiRequest(`/parts/search?q=${encodeURIComponent(query)}`);
  if (!response.ok) {
    showToast(await readErrorMessage(response));
    if (statusNode) statusNode.textContent = "Ошибка поиска";
    return;
  }
  const data = await response.json();
  const { items, hint } = normalizeListResponse(data);
  renderParts(items || []);
  if (statusNode) {
    statusNode.textContent = hint || `Найдено позиций: ${items.length}`;
  }
});

document.getElementById("parts-result")?.addEventListener("click", async (event) => {
  const button = event.target.closest("[data-add-part]");
  if (button) {
    if (!getUserToken()) {
      showToast("Войдите в аккаунт");
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
      showToast(await readErrorMessage(response));
      return;
    }
    showToast("Товар добавлен в корзину");
    loadCart();
    return;
  }

  const checkBtn = event.target.closest("[data-check-part]");
  if (!checkBtn) return;
  if (!getUserToken()) {
    showToast("Войдите в аккаунт");
    return;
  }

  const response = await apiRequest("/user/parts/check", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      part_id: checkBtn.dataset.checkPart || "",
      name: checkBtn.dataset.checkName || ""
    })
  });
  if (!response.ok) {
    showToast(await readErrorMessage(response));
    return;
  }
  const data = await response.json();
  showToast(data.message || "Уточнение выполнено");
});

document.getElementById("cart-container")?.addEventListener("click", async (event) => {
  const button = event.target.closest("[data-remove-part]");
  if (!button) return;
  const response = await apiRequest(`/user/cart/items/${button.dataset.removePart}`, { method: "DELETE" });
  if (!response.ok) {
    showToast("Не удалось удалить товар");
    return;
  }
  showToast("Позиция удалена");
  loadCart();
});

document.getElementById("checkout-form")?.addEventListener("submit", async (event) => {
  event.preventDefault();
  const statusNode = document.getElementById("checkout-status");
  const payload = { address: document.getElementById("checkout-address")?.value.trim() || "" };
  if (!payload.address) {
    if (statusNode) statusNode.textContent = "Введите адрес доставки.";
    return;
  }
  const response = await apiRequest("/user/cart/checkout", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });
  if (!response.ok) {
    if (statusNode) statusNode.textContent = await readErrorMessage(response);
    return;
  }
  const data = await response.json();
  const orderId = data.order_id ?? data.id;
  const amount = data.amount ?? data.total ?? 0;
  if (statusNode) {
    statusNode.textContent = orderId
      ? `Заказ #${orderId} создан.${amount ? ` Сумма: ${money(amount)}.` : ""} ${data.message || ""}`.trim()
      : (data.message || "Заказ создан.");
  }
  showToast("Заказ оформлен");
  if (data.payment_url) {
    window.open(data.payment_url, "_blank", "noopener");
  }
  loadCart();
  loadUserOrders();
});

document.getElementById("register-form")?.addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  clearFieldErrors(form);
  setFormError("register-error", "");
  const submitBtn = form.querySelector('button[type="submit"]');
  if (submitBtn) submitBtn.disabled = true;

  try {
    const payload = {
      name: document.getElementById("reg-name")?.value.trim() || "",
      phone: normalizePhone(document.getElementById("reg-phone")?.value || ""),
      email: document.getElementById("reg-email")?.value.trim() || "",
      password: document.getElementById("reg-password")?.value || "",
      address: document.getElementById("reg-address")?.value.trim() || ""
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
    if (String(payload.password).trim().length < 6) {
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
    saveProfile({
      name: payload.name,
      email: payload.email,
      phone: payload.phone,
      address: payload.address
    });
    setUserToken(tokens.access_token);
    showToast("Добро пожаловать!");
    afterAuthSuccess();
  } catch (err) {
    const msg = err && err.message ? err.message : String(err);
    setFormError("register-error", `Ошибка: ${msg}`);
    showToast(`Ошибка: ${msg}`);
  } finally {
    if (submitBtn) submitBtn.disabled = false;
  }
});

document.getElementById("login-form")?.addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  clearFieldErrors(form);
  setFormError("login-error", "");
  const submitBtn = form.querySelector('button[type="submit"]');
  if (submitBtn) submitBtn.disabled = true;

  try {
    const payload = {
      email: document.getElementById("login-email")?.value.trim() || "",
      password: document.getElementById("login-password")?.value || ""
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
    saveProfile({ email: payload.email });
    setUserToken(data.access_token);
    showToast("Вы вошли в аккаунт");
    afterAuthSuccess();
  } catch (err) {
    const msg = err && err.message ? err.message : String(err);
    setFormError("login-error", `Ошибка: ${msg}`);
    showToast(`Ошибка: ${msg}`);
  } finally {
    if (submitBtn) submitBtn.disabled = false;
  }
});

document.getElementById("car-form")?.addEventListener("submit", async (event) => {
  event.preventDefault();
  const payload = {
    name: document.getElementById("car-name")?.value.trim() || "",
    vin: (document.getElementById("car-vin")?.value || "").trim().toUpperCase()
  };
  const response = await apiRequest("/user/cars", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });
  if (!response.ok) {
    showToast(await readErrorMessage(response));
    return;
  }
  showToast("Авто добавлено");
  document.getElementById("car-form")?.reset();
  loadCars();
});

document.getElementById("cars-list")?.addEventListener("click", async (event) => {
  const button = event.target.closest("[data-car-delete]");
  if (!button) return;
  const response = await apiRequest(`/user/cars/${button.dataset.carDelete}`, { method: "DELETE" });
  if (!response.ok) {
    showToast("Ошибка удаления авто");
    return;
  }
  showToast("Авто удалено");
  loadCars();
});

updateSessionLabel();
applyRoute();
if (getUserToken()) {
  afterAuthSuccess();
}
