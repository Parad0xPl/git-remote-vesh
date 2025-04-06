# Vesh
Vesh is a simple tool designed to manage remote Git repositories stored inside **VeraCrypt** containers. It provides seamless integration with Git by implementing a custom Git remote helper (`git-remote-vesh`) and supports remote access via **SSHFS**.

## Features
- **Encrypted Repository**: Securely store Git repositories in VeraCrypt containers.
- **Remote Access**: Use SSHFS to mount remote directories and store them offsite.
- **Cross-Platform Support**: Works on both Windows and Linux platforms.
- **Custom Git Remote Helper**: Simplifies Git operations (`push`, `pull`, `fetch`) with encrypted repositories.

## Dependencies
### Windows
- [SSHFS-Win](https://github.com/winfsp/sshfs-win)
- [VeraCrypt](https://www.veracrypt.fr/en/Home.html)

### Linux
- [sshfs](https://github.com/libfuse/sshfs)
- [VeraCrypt](https://www.veracrypt.fr/en/Home.html)

## Installation
1. Install the required dependencies for your platform:
   - On **Windows**, install [SSHFS-Win](https://github.com/winfsp/sshfs-win) and [VeraCrypt](https://www.veracrypt.fr/en/Home.html).
   - On **Linux**, install [sshfs](https://github.com/libfuse/sshfs) and [VeraCrypt](https://www.veracrypt.fr/en/Home.html).
2. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/vesh.git
   ```
3. Add the `git-remote-vesh` script to your system's PATH to make it accessible globally.

## How It Works
1. **Remote URL Parsing**:
   - The `vesh` protocol is used in the remote URL (e.g., `vesh://user@host:/path/to/repo.git`).
   - The tool extracts the SSH credentials and the path to the VeraCrypt container.

2. **Mounting**:
   - The remote directory is mounted using **SSHFS**.
   - The VeraCrypt container is unlocked and mounted locally to access the Git repository.

3. **Git Operations**:
   - Once the repository is accessible, standard Git commands like `push`, `pull`, and `fetch` are executed on the mounted repository.

4. **Cleanup**:
   - After the Git operation is complete, the tool automatically unmounts the SSHFS directory and the VeraCrypt container to ensure security.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

