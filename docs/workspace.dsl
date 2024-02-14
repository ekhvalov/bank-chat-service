# https://docs.structurizr.com/dsl/basics#basics
# To see diagrams:
#   1. Run dependencies COMPOSE_PROFILES=structurizr-ui task deps
#   2. Open http://localhost:8070
workspace {
    name "Bank Chat System"
    description "System that allows clients to chat with support managers"

    model {
        client = person "Client" "Bank client who wants to solve a problem"
        manager = person "Manager" "Bank manager who can solve the problem"

        chatSystem = softwareSystem "Bank Chat System" "System that allows clients to chat with support managers" {
            clientWebApplication = container "Client UI" "Javascript" {
                tags "Web Browser"
            }
            managerWebApplication = container "Manager UI" "Javascript" {
                tags "Web Browser"
            }

            chat = container "Chat" "Golang" {
                # Manager components
                # Manager handlers
                managerHandlerSendMessage = component "Handler\n SendMessage" "Golang"
                managerHandlerGetChats = component "Handler\n GetChats" "Golang"
                managerHandlerFreeHands = component "Handler\n FreeHands" "Golang"
                managerHandlerFreeHandsAvailability = component "Handler\n FreeHandsAvailability" "Golang"
                managerHandlerGetChatHistory = component "Handler\n GetChatHistory" "Golang"
                managerHandlerCloseChat = component "Handler\n CloseChat" "Golang"

                # Manager use cases
                managerUseCaseSendMessage = component "Use Case\n SendMessage" "Golang"
                managerUseCaseGetChats = component "Use Case\n GetChats" "Golang"
                managerUseCaseFreeHands = component "Use Case\n FreeHands" "Golang"
                managerUseCaseCanReceiveProblem = component "Use Case\n CanReceiveProblem" "Golang"
                managerUseCaseGetChatHistory = component "Use Case\n GetChatHistory" "Golang"
                managerUseCaseCloseChat = component "Use Case\n CloseChat" "Golang"

                # Client components

                # Repositories
                problemsRepository = component "Problems Repository" "Entgo"
                messagesRepository = component "Messages Repository" "Entgo"
                chatsRepository = component "Chats Repository" "Entgo"
                jobsRepository = component "Jobs Repository" "Entgo"

                # Services
                managerPoolService = component "Manager Pool Service" "In-memory"
                managerLoadService = component "Manager Load Service" "?"
                afcProcessorService = component "AFC Processor Service" "Golang"
                outboxService = component "Outbox Service" "Golang"
            }


            chatDatabase = container "Database" "PostgreSQL" {
                tags "Database"
            }
            chatQueue = container "Queue" "Kafka" {
                tags "Queue"
            }
        }
        IAMSystem = softwareSystem "IAM System" "System that allows clients, and managers to authenticate" {
            tags "External System"
        }
        FPSystem = softwareSystem "Fraud Prevention System" "Monitors the appearance of fraud in the clients' messages" {
            tags "External System"
        }

        # Relationships to/from context
        client -> chatSystem "Uses" "Sends, receives messages"
        manager -> chatSystem "Uses" "Receives, sends messages"
        chatSystem -> IAMSystem "Uses" "Authenticates"
        chatSystem -> FPSystem "Uses" "Checks for fraud in messages"

        # Relationships to/from containers
        client -> clientWebApplication "Uses" "Sends, receives messages"
        manager -> managerWebApplication "Uses" "Receives, sends messages"
        managerWebApplication -> chat "Uses" "Sends, receives messages"
        clientWebApplication -> chat "Uses" "Sends, receives messages"
        chat -> IAMSystem "Uses" "Authenticates"
        FPSystem -> chatQueue "Uses" "Checks for fraud in messages"
        chat -> chatDatabase "Uses" "Stores messages"
        chat -> chatQueue "Uses" "Sends messages"

        # Relationships to/from components
        # Manager send message
        managerWebApplication -> managerHandlerSendMessage "HTTP" "POST"
        managerHandlerSendMessage -> managerUseCaseSendMessage "Uses" "Uses"
        managerUseCaseSendMessage -> messagesRepository "Creates"
        managerUseCaseSendMessage -> problemsRepository "Get assigned problem"
        manageruseCaseSendMessage -> outboxService "Put job"
        # Manager get chats
        managerWebApplication -> managerHandlerGetChats "HTTP" "GET"
        managerHandlerGetChats -> managerUseCaseGetChats "Uses" "Uses"
        managerUseCaseGetChats -> chatsRepository "Uses" "Reads"
        # Manager free hands
        managerWebApplication -> managerHandlerFreeHands "HTTP" "POST"
        managerHandlerFreeHands -> managerUseCaseFreeHands "Uses"
        managerUseCaseFreeHands -> managerPoolService "Contains"
        managerUseCaseFreeHands -> managerLoadService "CanManagerTakeProblem"
        # Manager can take problem
        managerWebApplication -> managerHandlerFreeHandsAvailability "HTTP" "POST"
        managerHandlerFreeHandsAvailability -> managerUseCaseCanReceiveProblem "Uses"
        managerUseCaseCanReceiveProblem -> managerPoolService "Contains"
        managerUseCaseCanReceiveProblem -> managerLoadService "CanManagerTakeProblem"
        # Manager close chat
        managerWebApplication -> managerHandlerCloseChat "HTTP" "POST"
        managerHandlerCloseChat -> managerUseCaseCloseChat "Uses"
        managerUseCaseCloseChat -> chatsRepository "Uses" "Get cient ID by chat ID"
        managerUseCaseCloseChat -> problemsRepository "Uses" "Resolve assigned problem
        managerUseCaseCloseChat -> messagesRepository "Uses" "Create service message"
        managerUseCaseCloseChat -> outboxService "Put job"
        # Manager get chat history
        managerWebApplication -> managerHandlerGetChatHistory "HTTP" "GET"
        managerHandlerGetChatHistory -> managerUseCaseGetChatHistory "Uses"
        managerUseCaseGetChatHistory -> problemsRepository "Uses" "Get unresolved problem"
        managerUseCaseGetChatHistory -> messagesRepository "Uses" "Get problem messages"

        # Repositories
        problemsRepository -> chatDatabase "Uses" "Reads, writes"
        chatsRepository -> chatDatabase "Uses" "Reads, writes"
        messagesRepository -> chatDatabase "Uses" "Reads, writes"
        jobsRepository -> chatDatabase "Uses" "Reads, writes"

        # Services
        outboxService -> jobsRepository "Uses" "Reads, writes"
        afcProcessorService -> messagesRepository "Visible for manager\n Block"
        afcProcessorService -> outboxService "Put job"
        afcProcessorService -> chatQueue "Uses" "Reads, writes"
    }

    views {
        theme https://static.structurizr.com/themes/default/theme.json

        styles {
            element "External System" {
                background #999999
                color #ffffff
            }
            element "Web Browser" {
                shape WebBrowser
            }
            element "Database" {
                shape Cylinder
            }
            element "Queue" {
                shape Pipe
            }
        }

        systemContext chatSystem "Context" {
            include *
        }

        container chatSystem "Container" {
            include *
        }

        component chat "Component" {
            include *
        }
    }

    configuration {
        scope softwaresystem
    }

}