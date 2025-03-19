flowchart TD
    A[Main Function] --> B[Load Config]
    B --> C[Migrate and Seed Database]
    C --> D[Initialize Worker Pool]
    D --> E[Start Worker Pool]
    E --> F[Create OS Signal Channel]
    F --> G[Run Providers in Parallel]
    G --> H[Fetch and Process Tasks]
    H --> I[Submit Tasks to Worker Pool]
    I --> J[Worker Pool Processes Tasks]
    J --> K[Wait for All Providers to Finish]
    K --> L[Stop Worker Pool]
    L --> M[Exit Application]

    subgraph Worker Pool
        E --> N[Start Workers]
        N --> O[Worker Goroutines]
        O --> P[Process Task Queue]
        P --> Q[Create Task in Database]
        Q --> R[Log Task Creation]
        R --> S[Wait for Context Done]
        S --> T[Stop Worker]
    end

    subgraph Providers
        G --> U[Provider 1]
        G --> V[Provider 2]
        U --> H
        V --> H
    end

    subgraph Channels and WaitGroups
        F --> W[OS Signal Channel]
        G --> X[WaitGroup for Providers]
        N --> Y[WaitGroup for Workers]
    end

    subgraph Task Processing
        H --> Z[Fetch Data from Provider]
        Z --> AA[Parse JSON]
        AA --> AB[Map to Task]
        AB --> AC[Submit Task to Worker Pool]
    end

    subgraph Worker Goroutine
        O --> AD[Receive Task from Task Queue]
        AD --> AE[Create Task in Database]
        AE --> AF[Log Task Creation]
        AF --> AG[Wait for Context Done]
        AG --> AH[Stop Worker]
    end

    subgraph Graceful Shutdown
        W --> AI[Receive OS Signal]
        AI --> AJ[Cancel Context]
        AJ --> AK[Stop Worker Pool]
        AK --> AL[Wait for Workers to Finish]
        AL --> AM[Exit Application]
    end

    style A fill:#f9f,stroke:#333,stroke-width:4px
    style D fill:#bb,stroke:#333,stroke-width:4px
    style E fill:#bb,stroke:#333,stroke-width:4px
    style G fill:#bb,stroke:#333,stroke-width:4px
    style H fill:#bb,stroke:#333,stroke-width:4px
    style I fill:#bb,stroke:#333,stroke-width:4px
    style J fill:#bb,stroke:#333,stroke-width:4px
    style K fill:#bb,stroke:#333,stroke-width:4px
    style L fill:#bb,stroke:#333,stroke-width:4px
    style M fill:#bb,stroke:#333,stroke-width:4px
    style N fill:#bb,stroke:#333,stroke-width:4px
    style O fill:#bb,stroke:#333,stroke-width:4px
    style P fill:#bb,stroke:#333,stroke-width:4px
    style Q fill:#bb,stroke:#333,stroke-width:4px
    style R fill:#bb,stroke:#333,stroke-width:4px
    style S fill:#bb,stroke:#333,stroke-width:4px
    style T fill:#bb,stroke:#333,stroke-width:4px
    style U fill:#bb,stroke:#333,stroke-width:4px
    style V fill:#bb,stroke:#333,stroke-width:4px
    style W fill:#bb,stroke:#333,stroke-width:4px
    style X fill:#bb,stroke:#333,stroke-width:4px
    style Y fill:#bb,stroke:#333,stroke-width:4px
    style Z fill:#bb,stroke:#333,stroke-width:4px
    style AA fill:#bb,stroke:#333,stroke-width:4px
    style AB fill:#bb,stroke:#333,stroke-width:4px
    style AC fill:#bb,stroke:#333,stroke-width:4px
    style AD fill:#bb,stroke:#333,stroke-width:4px
    style AE fill:#bb,stroke:#333,stroke-width:4px
    style AF fill:#bb,stroke:#333,stroke-width:4px
    style AG fill:#bb,stroke:#333,stroke-width:4px
    style AH fill:#bb,stroke:#333,stroke-width:4px
    style AI fill:#bb,stroke:#333,stroke-width:4px
    style AJ fill:#bb,stroke:#333,stroke-width:4px
    style AK fill:#bb,stroke:#333,stroke-width:4px
    style AL fill:#bb,stroke:#333,stroke-width:4px
    style AM fill:#bb,stroke:#333,stroke-width:4px