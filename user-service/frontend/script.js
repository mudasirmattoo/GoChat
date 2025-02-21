document.addEventListener("DOMContentLoaded", () => {
  console.log("JS is loaded");
  
  const showLoginBtn = document.getElementById("showLogin");
  const showRegisterBtn = document.getElementById("showRegister");
  const loginForm = document.getElementById("loginForm");
  const registerForm = document.getElementById("registerForm");
  const messageDiv = document.getElementById("message");


  registerForm.classList.add("hidden");
  
  showLoginBtn.classList.add("bg-blue-500", "text-white");
  showLoginBtn.classList.remove("bg-gray-300");
  
  showRegisterBtn.classList.add("bg-gray-300");
  showRegisterBtn.classList.remove("bg-blue-500", "text-white");

  // Toggle Forms
  showLoginBtn.addEventListener("click", () => {
    loginForm.classList.remove("hidden");
    registerForm.classList.add("hidden");

    showLoginBtn.classList.add("bg-blue-500", "text-white");
    showLoginBtn.classList.remove("bg-gray-300");
    
    showRegisterBtn.classList.add("bg-gray-300");
    showRegisterBtn.classList.remove("bg-blue-500", "text-white");
  });

  showRegisterBtn.addEventListener("click", () => {
    registerForm.classList.remove("hidden");
    loginForm.classList.add("hidden"); 

    showRegisterBtn.classList.add("bg-blue-500", "text-white");
    showRegisterBtn.classList.remove("bg-gray-300");
    
    showLoginBtn.classList.add("bg-gray-300");
    showLoginBtn.classList.remove("bg-blue-500", "text-white");
  });

  function showMessage(message, isError = false) {
    messageDiv.textContent = message;
    
    if (isError) {
      messageDiv.className = "mt-4 text-center text-red-500 font-medium animate-pulse";
    } else {
      messageDiv.className = "mt-4 text-center text-green-500 font-medium";
    }
    
    setTimeout(() => {
      messageDiv.textContent = "";
      messageDiv.className = "mt-4 text-center";
    }, 5000);
  }


  registerForm.addEventListener("submit", async (event) => {
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
      const response = await fetch("http://localhost:9080/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, email, password, confirm_password: confirmPassword }),
        credentials: "include"  // For cookies if you're using session-based auth
      });

      const result = await response.text();
      
      if (response.ok) {
        showMessage("Registration successful!");
        registerForm.reset();

        setTimeout(() => showLoginBtn.click(), 1500);
      } else {
        showMessage(result || "Registration failed", true);
      }
    } catch (error) {
      console.error("Registration error:", error);
      showMessage("Error connecting to server. Please try again.", true);
    }
  });


  loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    
    const username = document.getElementById("loginUsername").value.trim();
    const password = document.getElementById("loginPassword").value;


    if (!username || !password) {
      showMessage("Username and password are required", true);
      return;
    }

    try {
      const response = await fetch("http://localhost:9080/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
        credentials: "include"  
      });

      if (response.ok) {
        const result = await response.text();
        showMessage("Login successful!");
        loginForm.reset();
        // Redirect to dashboard or home page after successful login
        window.location.href = "/dashboard.html";
      } else {
        const errorText = await response.text();
        showMessage(errorText || "Invalid username or password", true);
      }
    } catch (error) {
      console.error("Login error:", error);
      showMessage("Error connecting to server. Please try again.", true);
    }
  });

  console.log("JavaScript initialized successfully - toggle should be working!");
});

