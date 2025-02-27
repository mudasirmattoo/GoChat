  document.addEventListener("DOMContentLoaded", () => {
    // console.log("JS is loaded");
    
    const loginForm = document.getElementById("loginForm");
    const registerForm = document.getElementById("registerForm");
    const messageDiv = document.getElementById("message");
  
    function showMessage(message, isError = false) {
      messageDiv.textContent = message;
      
      if (isError) {
        messageDiv.className = "mt-4 text-center text-red-500 font-medium";
      } else {
        messageDiv.className = "mt-4 text-center text-green-500 font-medium";
      }
      
      setTimeout(() => {
        messageDiv.textContent = "";
        messageDiv.className = "mt-4 text-center";
      }, 5000);
    }
  
    function checkFlashMessages() {
      const flashCookie = getCookie("flash_message");
      if (flashCookie) {
        showMessage(decodeURIComponent(flashCookie));
        document.cookie = "flash_message=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      }
    }
  
    function getCookie(name) {
      const value = `; ${document.cookie}`;
      const parts = value.split(`; ${name}=`);
      if (parts.length === 2) return parts.pop().split(';').shift();
    }
  
    checkFlashMessages();
  
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
          body: JSON.stringify({ username, email, password, confirm_password: confirmPassword }),
          credentials: "include"  
        });
  
        if (response.ok) {
          const result = await response.json();
          
          showMessage(result.message || "Registration successful!");
          registerForm.reset();
          
          if (result.redirect_url) {
            setTimeout(() => {
              window.location.href = result.redirect_url;
            }, 1000); 
          } else {
            setTimeout(() => {
              window.location.href = "/login";
            }, 500);
          }
        } else {
          const errorText = await response.text();
          try {
            const errorJson = JSON.parse(errorText);
            showMessage(errorJson.message || "Registration failed", true);
          } catch (e) {
            showMessage(errorText || "Registration failed", true);
          }
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
          credentials: "include"  
        });
  
        if (response.ok) {
          const result = await response.json();
          
          showMessage(result.message || "Login successful!");
          loginForm.reset();
          
          if (result.redirect_url) {
            setTimeout(() => {
              window.location.href = result.redirect_url;
            }, 500); 
          } else {
            setTimeout(() => {
              window.location.href = "/dashboard";
            }, 1000);
          }
        } else {
          const errorText = await response.text();
          try {
            const errorJson = JSON.parse(errorText);
            showMessage(errorJson.message || "Invalid username or password", true);
          } catch (e) {
            showMessage(errorText || "Invalid username or password", true);
          }
        }
      } catch (error) {
        console.error("Login error:", error);
        showMessage("Error connecting to server. Please try again.", true);
      }
    });
  
    // console.log("JavaScript initialized successfully");
  });