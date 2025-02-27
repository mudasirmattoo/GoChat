  document.addEventListener("DOMContentLoaded", () => {
    
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
      }, 2000);
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
          body: JSON.stringify({ username, email, password, confirm_password: confirmPassword }),
          // redirect: "follow",
          credentials: "include"  
        });

        if(!response.ok){
          const errorText = await response.text();
          try{
            const errorJson = JSON.parse(errorText);
            showMessage(errorJson.message || "Registration failed ", true);
          } catch(e) {
            showMessage(errorText || "Registration failed", true);
          }
        }else{
          const data = await response.json();
          showMessage("Registration successful!");

          localStorage.setItem("token",data.token);
          window.location.href = "/login";

        }

 
        }catch(error){
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
          showMessage("Login successful!");
    
          localStorage.setItem("token", data.token);
    
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
    
  });