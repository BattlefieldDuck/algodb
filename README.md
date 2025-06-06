# algodb

A Go-based tool to find all possible solutions of different Rubik’s Cube states. It reads scramble sequences from CSV files, computes solutions using a solver, and saves the results back to CSV in a `db` folder. The process is automated via GitHub Actions.

Users can find the computed algorithms in the [db](db) folder, which is automatically generated by [GitHub Actions](https://github.com/BattlefieldDuck/algodb/actions).

```mermaid
---
config:
  theme: redux-dark-color
  layout: dagre
  look: classic
---
sequenceDiagram
    participant GA as GitHub Actions
    participant Prog as algodb Program
    participant CFG as Config CSV
    participant Solver as Cube Solver
    participant DB as DB Folder (CSV)
    participant GH as GitHub API

    GA->>Prog: Trigger job
    Prog->>CFG: Read scramble CSV files
    CFG-->>Prog: Return scramble data
    Prog->>Prog: Parse scramble moves
    Prog->>Solver: Compute solution
    Solver-->>Prog: Return solution data
    Prog->>Prog: Format solution data
    Prog->>DB: Write solution CSV files
    GA->>GH: Create Pull Request
    GH-->>GA: Pull Request Created
```

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.
