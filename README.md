<div align="center">
  <img src="https://github.com/user-attachments/assets/c531f18b-5886-464a-a189-971b39134aee" alt="Browsir Logo" width="250">
</div>

# Browsir 🎩

A simple yet powerful command-line tool to manage multiple browser profiles and shortcuts.

## Features ✨

- Launch different browser profiles with a single command
- Create and manage shortcuts to your favorite websites
- Support for both global (config file) and local shortcuts
- Smart shortcut suggestions when typos occur
- Interactive shortcut creation
- Support for Firefox, Chrome, Brave, Zen and Arc browsers (on macOS)
- Cross-platform: works on macOS, Linux and Windows

## Installation 🚀

### From Source

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/browsir.git
   cd browsir
   ```

2. Build and install:

   **READ ME Before you continue the installation**

If you already used `browsir`, version 2.0.0 has breaking changes for configuration files.

This will force you to enter all the previously saved informations again, so be sure to save your configurations somewhere else, clone and then swap them with the installed ones which are now in `$HOME/.config/browsir/(config.yml|links|shortcuts)`

```bash
make install
```

This will build the binary and install it to `$HOME/bin`

3. Verify the installation:
   ```bash
   browsir --version
   ```

You can also just build without installing:

```bash
make build
./browsir --version
```

## Update

```bash
make update
```

This will build the binary and install it to `$HOME/bin`

Verify the installation:
```bash
browsir --version
```

## Usage 📖

```bash
# Open browser with profile
browsir [profile] [shortcut | website]

browsir personal mail
browsir personal gmail.com

# Search on google, duckduckgo and bravesearch
# Default search engine is google
browsir [profile] [-se | --search-engine]=[google | brave | duckduckgo] -q=[your query]

browsir personal -q="What's the distance between the moon and the sun"
browsir personal -se=brave -q="Is Brave better for privacy"

# Manage links and shortcuts
browsir add link <link> -c="<categories>"    # Add a link with categories
browsir add shortcut <shortcut> <url>      # Add a local shortcut, do not include http:// or https://
browsir rm link <link>                     # Remove a link
browsir rm shortcut <shortcut>             # Remove a local shortcut
browsir list links                         # List all links
browsir list all                           # List all links and categories
browsir preview <link>                     # Preview a link
```

## Available commands and flags

```bash
-ls, --list-shortcuts, list all shortcuts
-se, --search-engine, set search engine for search
-q, query search engine
-v, --version, check browsir version
-h, --help, help
```

### Configuration 🔧

1. Create or modify `.browsir.yml` in your browsir directory:
   ```yaml
   app_name: browsir
   browser_name: chrome # can be 'firefox', 'chrome', 'brave', 'zen' or 'arc'
   profiles:
     - name: personal # profile name you'll use in commands
       profile_dir: Default # actual profile directory name
       description: Personal browsing
     - name: work
       profile_dir: Profile 1
       description: Work profile
   shortcuts:
     google: google.com
     github: github.com
     mail: gmail.com
   ```

The configuration file allows you to:

- Set your preferred browser (`firefox`, `chrome`, `brave`, `arc`, `zen`)
- Define multiple browser profiles with custom names
- Add global shortcuts to frequently visited websites

You can find your Chrome profile directory names by visiting:

- Chrome: `chrome://version`
- Brave: `brave://version`
- Arc: `arc://version`
- Firefox, Firefox Developer Edition & Zen: `about:profiles`

> [!NOTE]
> When using both Firefox, Firefox Developer Edition or Zen on MacOS, ensure each app uses its appropriate profile.
> For example, the "My Firefox Developer Edition Profile" in `about:profiles` should always be opened with the Firefox Developer Edition app.
> Using the wrong profile with either app will cause them to crash.

## Maintainers 👨‍💻👩‍💻

<div align="center">
  <table>
    <tr>
      <td align="center">
        <a href="https://github.com/404answernotfound">
          <img src="https://github.com/404answernotfound.png" width="100px;" alt="Lorenzo Pieri"/>
          <br />
          <sub>
            <b>Lorenzo Pieri</b>
          </sub>
        </a>
        <br />
        <span>💻 Maintainer</span>
      </td>
      <td align="center">
        <a href="https://github.com/Wabri">
          <img src="https://github.com/Wabri.png" width="100px;" alt="Gabriele Puliti"/>
          <br />
          <sub>
            <b>Gabriele Puliti</b>
          </sub>
        </a>
        <br />
        <span>💻 Maintainer</span>
      </td>
    </tr>
  </table>
</div>