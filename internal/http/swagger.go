package http

import (
	"go-api/internal"

	"fmt"
	"net/http"
	"os"
	"strconv"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Load Swagger Docs on
func (s *Server) loadSwagger() {
	sm := http.NewServeMux()

	// Load Custom CSS
	sm.HandleFunc("GET /swagger/custom.css", func(w http.ResponseWriter, r *http.Request) {
		// Open the CSS file
		file, err := os.Open("docs/swagger/custom.css")
		if err != nil {
			internal.APIError(w, "Swagger::loadSwagger", "Unable to load custom css", http.StatusInternalServerError, err)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			internal.APIError(w, "Swagger::loadSwagger", "Unable to load information of custom css", http.StatusInternalServerError, err)
			return
		}

		// Set the content type header to "text/css"
		w.Header().Set("Content-Type", "text/css")

		// Serve the file content using http.ServeContent
		http.ServeContent(w, r, "", fileInfo.ModTime(), file)
	})

	// Load /swagger & /swagger/doc.json
	sm.HandleFunc("GET /", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+strconv.Itoa(s.Port)+"/swagger/doc.json"),
		httpSwagger.AfterScript(fmt.Sprintf(`
			// Set Input Value Natively
			function setNativeValue(element, value) {
				const valueSetter = Object.getOwnPropertyDescriptor(element, 'value').set;
				const prototype = Object.getPrototypeOf(element);
				const prototypeValueSetter = Object.getOwnPropertyDescriptor(prototype, 'value').set;

				if (valueSetter && valueSetter !== prototypeValueSetter) {
					prototypeValueSetter.call(element, value);
				} else {
					valueSetter.call(element, value);
				}
			}

			// Function to trigger click on the button with classes "btn authorize unlocked"
			function triggerAuth(token) {
				// Select the button with the specified classes
				var authorizeButton = document.querySelector('.btn.authorize.unlocked');

				// Check if the button is found
				if (authorizeButton) {
					// Trigger the click event on the button
					authorizeButton.click();

					// Call the function to set input content after a delay
					// Set a timeout to fill the input field with token data
					setTimeout(function(token) {
						// Get the input field
						var inputField = document.querySelector('.auth-container input');

						// Fill the input field with token data
						setNativeValue(inputField, token);
						inputField.dispatchEvent(new Event('input', { bubbles: true }));
						document.querySelector('.btn.modal-btn.authorize').click();
						document.querySelector('.btn.modal-btn.btn-done').click();
					}, 100, token);
				} else {
					console.error('Authorize Button not found');
				}
			}

			function addCSS(filePath) {
				// Add custom style
				var link = document.createElement('link');
        link.rel = 'stylesheet';
				link.type = 'text/css';
				link.href = filePath;
				document.head.appendChild(link);
			}
			`),
		),
		httpSwagger.UIConfig(map[string]string{
			// responseInterceptor to process signin response upon success
			"responseInterceptor": `(res) => {
				// Check if signin is successful
				if(res.url.includes("api/v1/auth/signin") && res.status == 200 && res.body.token) {
					// Trigger Auth Popup
					// https://stackoverflow.com/questions/76654069/automatically-authorize-swagger-endpoints-after-login-with-token
					triggerAuth("Bearer " + res.body.token);

					// Non Working flows:
					// window.ui.preauthorizeApiKey("Authorization", "Bearer " + res.body.token);
					// window.ui.clientAuthorizations.add('apiKey', new SwaggerClient.ApiKeyAuthorization('Authorization', 'Bearer ' + res.body.token, 'header'));
				}
			}`,
			"onComplete": fmt.Sprintf(
				`() => {
          addCSS("%s");
          addCSS("%s");
        }`,
				"http://localhost:"+strconv.Itoa(s.Port)+"/swagger/custom.css",
				"https://fonts.googleapis.com/css2?family=Lato:ital,wght@0,100;0,300;0,400;0,700;0,900;1,100;1,300;1,400;1,700;1,900&display=swap",
			),
		}),
	))
	s.Router.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swagger-ui"))))
	s.Router.Handle("/swagger/*", sm)

	internal.Debug("Http::Swagger", "         Docs: http://localhost:"+strconv.Itoa(s.Port)+"/swagger")
}
