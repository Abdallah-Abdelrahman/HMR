# HMR (Hot Module Replacement)

HMR is a lightweight tool that updates your browser in real-time without requiring a page refresh.

<a href="https://ibb.co/4Rf7fMXW"><img width='100%' src="https://i.ibb.co/GQJ0Jk87/output.gif" alt="output" border="0"></a>

## Motivation

When working on projects that do not use UI libraries like React or Vue—which come with built-in Hot Module Replacement—developers often miss out on the benefits of live updates. While tools like [browser-sync](https://www.npmjs.com/package/browser-sync) offer similar functionality, they can be resource-intensive and heavy on system memory. HMR was created as a lean alternative to speed up development by instantly reflecting UI changes without the overhead.

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Abdallah-Abdelrahman/HMR.git
   ```
2. **Navigate to the project directory:**
   ```bash
   cd HMR
   ```
3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

## Usage

1. **Build the static binary:**
   ```bash
   make build-static
   ```
2. **Create a symlink to the binary:**
   ```bash
   sudo ln -s <project-absolute-path>/bin/hmr /usr/local/bin/hmr
   ```
3. **Run HMR from anywhere:**
   ```bash
   hmr <path-to-directory-containing-html|css|js-files>
   ```
4. **Embed the client-side script:**  
   Include this [script](ws.html) in your HTML file to enable live updates.
