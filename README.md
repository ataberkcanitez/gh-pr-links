# pr-links
### GitHub CLI Pull Request Viewer

pr-links is a command-line tool written in Go that enhances the functionality of the GitHub CLI by displaying Pull Request information along with their URLs. The GitHub CLI itself lacks the ability to show the Pull Request URL, and this tool bridges that gap.

### Installation

Since it's GitHub CLI Extension, make sure you have [GitHub CLI](https://cli.github.com) installed on your system.

To install the extension, please execute the following command:

```bash
$ gh extension install ataberkcanitez/gh-pr-links
```

### Usage

After installing the extension, please execute following command to use:

```bash
$ gh pr-links
```

### Options

* `-h --help`: Show help for command.
* `-s --style <string>`: Sets the style of the output. Possible values include `StyleCompactLite`, `StyleUnicode`, `StyleDefault`, `StyleCompact`, `StyleMarkdown`, `StyleRounded`, and `StyleCompactClassic`.
* `-u --use-emoji <bool>`: Use emoji in the output. Possible values are `true` or `false`.
* `-o --org <string>`: Only show pull requests from the given organization. Separate multiple organizations with `|` (e.g. `--org="acme|globex"`), or use `-` to remove the filter.

Your choices are saved to `~/gh-pr-links-configuration.json`, so you only need to set them once.

Options Example:
```bash
$ gh pr-links --style=StyleMarkdown --use-emoji=true
```

### Example Output

The tool fetches open Pull Requests that are currently assigned to you and presents them in a formatted table. 
The table includes columns for the repository name, title, author, and a direct link to the Pull Request on GitHub.
Here is some example outputs for different styles:

![example-output](assets/outputExamples/rounded%20-%20no-emoji.png)
![example-output](assets/outputExamples/rounded%20-%20with-emoji.png)

![example-output](assets/outputExamples/compactLite%20-%20no-emoji.png)
![example-output](assets/outputExamples/compactLite%20-%20with-emoji.png)

![example-output](assets/outputExamples/compact%20-%20no-emoji.png)
![example-output](assets/outputExamples/compact%20-%20with-emoji.png)

![example-output](assets/outputExamples/compactClassic%20-%20no%20emoji.png)
![example-output](assets/outputExamples/compactClassic%20-%20with-emoji.png)

![example-output](assets/outputExamples/markdown%20-%20no-emoji.png)
![example-output](assets/outputExamples/markdown%20-%20with-emoji.png)

![example-output](assets/outputExamples/unicode%20-%20no-emoji.png)
![example-output](assets/outputExamples/unicode%20-%20with-emoji.png)

![example-output](assets/outputExamples/default%20-%20no-emoji.png)
![example-output](assets/outputExamples/default%20-%20with-emoji.png)

### License
This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.


### Acknowledgments

- GitHub CLI for providing the foundation for interacting with GitHub from the command line.
- alexeyco/simpletable for simplifying the creation of formatted tables in the console.

### Development

```bash
go test ./...            # run the tests
go build -o gh-pr-links . && gh extension install .   # try your local build
```

The code is organised as:

* `main.go` – wires everything together.
* `internal/cli` – flag parsing and the command flow.
* `internal/config` – reading and writing the configuration file.
* `internal/github` – fetching pull requests via `gh` and filtering them.
* `internal/render` – drawing the table.

### Contributing

Feel free to contribute by opening issues or creating pull requests. Your feedback and involvement are highly encouraged!

---
Enjoy using pr-links and feel free to reach out with any feedback or suggestions!