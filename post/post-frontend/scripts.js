const API_BASE_URL = "http://localhost:9081/posts";


async function fetchPosts() {
    try {
        const response = await fetch(API_BASE_URL);
        const posts = await response.json();

        const postsContainer = document.getElementById("posts-container");
        postsContainer.innerHTML = ""; // Clear previous content

        posts.forEach(post => {
            const postElement = document.createElement("div");
            postElement.classList.add("bg-white", "shadow-md", "p-4", "rounded-lg");

            postElement.innerHTML = `
                <h3 class="text-lg font-bold">${post.title}</h3>
                <p class="text-gray-700">${post.content}</p>
                <div class="flex gap-2 mt-2">
                    <button onclick="editPost(${post.id})" class="bg-blue-500 text-white px-4 py-2 rounded">Edit</button>
                    <button onclick="deletePost(${post.id})" class="bg-red-500 text-white px-4 py-2 rounded">Delete</button>
                </div>
            `;
            postsContainer.appendChild(postElement);
        });
    } catch (error) {
        console.error("Error fetching posts:", error);
    }
}


async function createPost(event) {
    event.preventDefault();

    const title = document.getElementById("title").value;
    const content = document.getElementById("content").value;

    try {
        const response = await fetch(API_BASE_URL, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title, content })
        });

        if (response.ok) {
            fetchPosts();
            document.getElementById("post-form").reset();
        } else {
            console.error("Failed to create post");
        }
    } catch (error) {
        console.error("Error creating post:", error);
    }
}


async function deletePost(postId) {
    try {
        const response = await fetch(`${API_BASE_URL}/${postId}`, {
            method: "DELETE",
        });

        if (response.ok) {
            fetchPosts();
        } else {
            console.error("Failed to delete post");
        }
    } catch (error) {
        console.error("Error deleting post:", error);
    }
}


async function editPost(postId) {
    try {
        const response = await fetch(`${API_BASE_URL}/${postId}`);
        const post = await response.json();

        document.getElementById("title").value = post.title;
        document.getElementById("content").value = post.content;
        document.getElementById("post-id").value = post.id;
    } catch (error) {
        console.error("Error fetching post:", error);
    }
}

async function updatePost(event) {
    event.preventDefault();

    const postId = document.getElementById("post-id").value;
    const title = document.getElementById("title").value;
    const content = document.getElementById("content").value;

    try {
        const response = await fetch(`${API_BASE_URL}/${postId}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title, content })
        });

        if (response.ok) {
            fetchPosts();
            document.getElementById("post-form").reset();
            document.getElementById("post-id").value = "";
        } else {
            console.error("Failed to update post");
        }
    } catch (error) {
        console.error("Error updating post:", error);
    }
}


document.getElementById("post-form").addEventListener("submit", (event) => {
    const postId = document.getElementById("post-id").value;
    if (postId) {
        updatePost(event);
    } else {
        createPost(event);
    }
});


document.addEventListener("DOMContentLoaded", fetchPosts);
