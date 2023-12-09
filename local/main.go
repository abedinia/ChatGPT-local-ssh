package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func runSSHCommand(message string) (string, error) {
	privateKeyPath := "/root/.ssh/id_rsa"
	cmd := exec.Command("ssh", "-i", privateKeyPath, "-o", "StrictHostKeyChecking=no", "USER@IP", fmt.Sprintf("/root/chatgpt/main '%s'", message))

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing SSH command: %v, Output: %s\n", err, output)
	}
	return string(output), err
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get the message from the form
	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Run the SSH command with the message
	response, err := runSSHCommand(message)
	if err != nil {
		// Logging the error to the server console
		log.Printf("Error running SSH command: %v", err)
		http.Error(w, fmt.Sprintf("Error executing SSH command: %v", err), http.StatusInternalServerError)
		return
	}

	// Send the SSH command response back to the client
	fmt.Fprint(w, response)
}

func serveForm(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
<html>
<head>
    <title>SSH Command Interface</title>
    <style>
        /* Basic styling for the webpage */
        body {
            font-family: Arial, sans-serif;
            background-color: #f0f0f0;
            margin: 0;
            height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
        }
        .container {
            background-color: #e8f5e9; /* Light green background */
            width: 800px;
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        input[type=text], input[type=submit] {
            width: 100%;
            padding: 10px;
            margin: 10px 0;
            border: 1px solid #ddd;
            border-radius: 5px;
            box-sizing: border-box;
        }
        input[type=submit] {
            background-color: #4CAF50; /* Green button */
            color: white;
            border: none;
            cursor: pointer;
        }
        input[type=submit]:hover {
            background-color: #45a049;
        }
        #response, #loading {
            color: green;
            margin-top: 20px;
        }
        #loading {
            display: none;
        }
        textarea {
            width: 100%; /* Full width */
            padding: 10px;
            margin: 10px 0;
            border: 1px solid #ddd;
            border-radius: 5px;
            box-sizing: border-box;
            resize: vertical; /* Allow vertical resizing */
        }
        
        #response {
            background-color: white; /* White background */
            border: 1px solid #ddd; /* Light grey border */
            padding: 10px;
            margin-top: 20px;
            height: 150px; /* Fixed height */
            overflow-y: auto; /* Enable vertical scrolling */
        }
    </style>
</head>
<body>
    <div class="container">
        <h1 style="color: green;">ChatGPT 3.5 Turbo</h1>
        <form id="sshForm">
            <textarea id="message" name="message" placeholder="Enter your message" rows="4" cols="50"></textarea>
            <input type="submit" value="Send">
        </form>
        <div id="response"></div>
        <div id="loading">Loading...</div>
    </div>

    <script>
        // JavaScript for handling form submission
        document.getElementById('sshForm').onsubmit = function(event) {
            event.preventDefault(); // Prevent the default form submission behavior
            document.getElementById('response').innerText = '';
            document.getElementById('loading').style.display = 'block'; // Show loading text

            var message = document.getElementById('message').value;
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/', true);
            xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
            xhr.onreadystatechange = function() {
                if (xhr.readyState == 4) {
                    document.getElementById('loading').style.display = 'none'; // Hide loading text
                    if (xhr.status == 200) {
                        document.getElementById('response').innerText = 'Response: ' + xhr.responseText;
                    } else {
                        document.getElementById('response').innerText = 'Error: ' + xhr.statusText;
                    }
                }
            };
            xhr.send('message=' + encodeURIComponent(message));
        };
        document.getElementById('message').addEventListener('keydown', function(event) {
            if (event.key === 'Enter' && !event.shiftKey) {
                event.preventDefault(); // Prevent the default action to avoid a line break in the textarea
                document.getElementById('sshForm').dispatchEvent(new Event('submit'));
            }
        });
    </script>
</body>
</html>
	`
	fmt.Fprint(w, html)
}

func main() {
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/form", serveForm)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
