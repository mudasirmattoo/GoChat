document.addEventListener("DOMContentLoaded", () => {
  const loginForm = document.getElementById("loginForm");
  const registerForm = document.getElementById("registerForm");
  const messageDiv = document.getElementById("message");

  function showMessage(message, isError = false) {
    if (!messageDiv) {
      return;
    }
    
    messageDiv.textContent = message;
  
    if (isError) {
      messageDiv.className = "mt-4 text-center text-red-500 font-medium";
    } else {
      messageDiv.className = "mt-4 text-center text-green-500 font-medium";
    }
  
    setTimeout(() => {
      if (messageDiv) {
        messageDiv.textContent = "";
        messageDiv.className = "mt-4 text-center";
      }
    }, 500);
  }

  function getAuthHeaders() {
    const token = localStorage.getItem("token");
    return token
      ? {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        }
      : {
          "Content-Type": "application/json",
        };
  }

  registerForm?.addEventListener("submit", async (event) => {
    event.preventDefault();

    const username = document.getElementById("registerUsername").value.trim();
    const email = document.getElementById("registerEmail").value.trim();
    const password = document.getElementById("registerPassword").value;
    const confirmPassword = document.getElementById("confirmPassword").value;

    if (!username || !email || !password) {
      showMessage("All fields are required", true);
      return;
    }

    if (password !== confirmPassword) {
      showMessage("Passwords do not match", true);
      return;
    }

    if (password.length < 6) {
      showMessage("Password must be at least 6 characters", true);
      return;
    }

    try {
      const response = await fetch("/register-user", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username,
          email,
          password,
          confirm_password: confirmPassword,
        }),
        credentials: "include",
      });

      if (!response.ok) {
        const errorText = await response.text();
        try {
          const errorJson = JSON.parse(errorText);
          showMessage(errorJson.message || "Registration failed ", true);
        } catch (e) {
          showMessage(errorText || "Registration failed", true);
        }
      } else {
        const data = await response.json();
        showMessage("Registration successful!");

        localStorage.setItem("token", data.token);
        updateUI();
        window.location.href = "/login";
      }
    } catch (error) {
      console.error("Registration error:", error);
      showMessage("Error connecting to server", true);
    }
  });

  loginForm?.addEventListener("submit", async (event) => {
    event.preventDefault();

    const username = document.getElementById("loginUsername").value.trim();
    const password = document.getElementById("loginPassword").value;

    if (!username || !password) {
      showMessage("Username and password are required", true);
      return;
    }

    try {
      const response = await fetch("/login-user", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
        credentials: "include",
      });

      if (response.ok) {
        const data = await response.json();
        // showMessage("Login successful!");
        localStorage.setItem("token", data.token);

        updateUI();
        window.location.href = "/dashboard";
      } else {
        const errorText = await response.text();
        try {
          const errorJson = JSON.parse(errorText);
          showMessage(errorJson.message || "Login failed", true);
        } catch (e) {
          showMessage(errorText || "Login failed", true);
        }
      }
    } catch (error) {
      console.error("Login error:", error);
      showMessage("Error connecting to server. Please try again.", true);
    }
  });

  const protectedPages = ["/dashboard"];
  const currentPath = window.location.pathname;

  if (protectedPages.includes(currentPath)) {
    const token = localStorage.getItem("token");
    if (!token) {
      window.location.href = "/login";
    }
  }

  const loginLink = document.getElementById("loginLink");
  const registerLink = document.getElementById("registerLink");
  const logoutButton = document.getElementById("logoutButton");

  function updateUI() {  
    const token = localStorage.getItem("token");  
    const isAuthenticated = !!token;  

    if (isAuthenticated) {  
        console.log("authenticated")
        if (loginLink) loginLink.classList.add("hidden");  
        if (registerLink) registerLink.classList.add("hidden");  
        if (logoutButton) {  
            logoutButton.classList.remove("hidden");  
            logoutButton.addEventListener("click", handleLogout);  
        }  
    } else {  
        if (loginLink) loginLink.classList.remove("hidden");  
        if (registerLink) registerLink.classList.remove("hidden");  
        if (logoutButton) logoutButton.classList.add("hidden");  
    }  
}  

  async function handleLogout() {
    try {
      localStorage.removeItem("token");

      document.cookie =
        "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";

      await fetch("/logout", {
        method: "POST",
        headers: getAuthHeaders(),
        credentials: "include",
      });

      showMessage("Logged out successfully!",true);

      updateUI();

      setTimeout(() => {
        window.location.href = "/";
      }, 1000);
    } catch (error) {
      console.error("Logout error:", error);
      showMessage("Error logging out", true);
    }
  }

  updateUI();
});
