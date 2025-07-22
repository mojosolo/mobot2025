// Package catalog provides agent communication protocol for multi-agent coordination
package catalog

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// AgentCommunicationSystem manages inter-agent messaging and state synchronization
type AgentCommunicationSystem struct {
	database       *Database
	messageQueue   *MessageQueue
	stateManager   *StateManager
	eventBus       *EventBus
	agents         map[string]Agent
	subscriptions  map[string][]string // agent -> message types
	middleware     []Middleware
	metrics        *CommunicationMetrics
	mu             sync.RWMutex
}

// Agent represents a communication-enabled agent
type Agent interface {
	GetID() string
	GetType() string
	HandleMessage(message *Message) error
	GetState() AgentState
	SetState(state AgentState) error
	IsHealthy() bool
}

// Message represents inter-agent communication
type Message struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	From        string                 `json:"from"`        // sender agent ID
	To          string                 `json:"to"`          // recipient agent ID or "broadcast"
	Subject     string                 `json:"subject"`     // message topic
	Payload     map[string]interface{} `json:"payload"`
	Metadata    MessageMetadata        `json:"metadata"`
	Status      string                 `json:"status"`      // pending, delivered, failed, expired
	RetryCount  int                    `json:"retry_count"`
	MaxRetries  int                    `json:"max_retries"`
	CreatedAt   time.Time              `json:"created_at"`
	DeliveredAt *time.Time             `json:"delivered_at,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
}

// MessageMetadata contains message routing and processing information
type MessageMetadata struct {
	Priority    string            `json:"priority"`    // critical, high, medium, low
	Urgent      bool              `json:"urgent"`      // bypass normal queue
	Persistent  bool              `json:"persistent"`  // store in database
	Encrypted   bool              `json:"encrypted"`   // encrypt payload
	Compressed  bool              `json:"compressed"`  // compress payload
	Tags        []string          `json:"tags"`        // message categorization
	Headers     map[string]string `json:"headers"`     // custom headers
	TTL         time.Duration     `json:"ttl"`         // time to live
	Correlation string            `json:"correlation"` // correlation ID for request/response
}

// AgentState represents the current state of an agent
type AgentState struct {
	AgentID      string                 `json:"agent_id"`
	Status       string                 `json:"status"`       // idle, busy, error, maintenance
	CurrentTask  string                 `json:"current_task"` // task being processed
	Progress     float64                `json:"progress"`     // 0.0-1.0
	Capabilities []string               `json:"capabilities"` // what the agent can do
	Config       map[string]interface{} `json:"config"`       // agent configuration
	Metrics      AgentMetrics           `json:"metrics"`      // performance metrics
	LastActivity time.Time              `json:"last_activity"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// AgentMetrics tracks agent performance
type AgentMetrics struct {
	MessagesReceived    int           `json:"messages_received"`
	MessagesProcessed   int           `json:"messages_processed"`
	MessagesFailed      int           `json:"messages_failed"`
	AvgProcessingTime   time.Duration `json:"avg_processing_time"`
	LastProcessingTime  time.Duration `json:"last_processing_time"`
	ErrorRate           float64       `json:"error_rate"`
	ThroughputPerMinute float64       `json:"throughput_per_minute"`
	UptimePercentage    float64       `json:"uptime_percentage"`
}

// MessageQueue manages message delivery and queuing
type MessageQueue struct {
	queues     map[string]*Queue // priority-based queues
	processing map[string]*Message
	workers    int
	mu         sync.RWMutex
}

// Queue represents a priority-based message queue
type Queue struct {
	Messages   []*Message
	Priority   string
	MaxSize    int
	mu         sync.RWMutex
}

// StateManager manages agent state synchronization
type StateManager struct {
	states     map[string]*AgentState // agent_id -> state
	snapshots  []*StateSnapshot       // state history
	mu         sync.RWMutex
}

// StateSnapshot captures system state at a point in time
type StateSnapshot struct {
	ID        string                    `json:"id"`
	States    map[string]*AgentState    `json:"states"`
	Timestamp time.Time                 `json:"timestamp"`
	Metadata  map[string]interface{}    `json:"metadata"`
}

// EventBus manages event publishing and subscription
type EventBus struct {
	subscribers map[string][]EventHandler // event_type -> handlers
	events      []*Event                  // event history
	mu          sync.RWMutex
}

// Event represents a system event
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`    // agent that generated event
	Target    string                 `json:"target"`    // specific target or "all"
	Data      map[string]interface{} `json:"data"`
	Level     string                 `json:"level"`     // debug, info, warning, error, critical
	Tags      []string               `json:"tags"`
	Timestamp time.Time              `json:"timestamp"`
}

// EventHandler defines event processing function
type EventHandler func(event *Event) error

// Middleware defines message processing middleware
type Middleware func(message *Message, next func(*Message) error) error

// CommunicationMetrics tracks system-wide communication metrics
type CommunicationMetrics struct {
	TotalMessages       int           `json:"total_messages"`
	MessagesPerSecond   float64       `json:"messages_per_second"`
	AverageLatency      time.Duration `json:"average_latency"`
	QueueSizes          map[string]int `json:"queue_sizes"`
	ActiveAgents        int           `json:"active_agents"`
	HealthyAgents       int           `json:"healthy_agents"`
	ErrorRate           float64       `json:"error_rate"`
	ThroughputTrend     []float64     `json:"throughput_trend"`
	LastUpdated         time.Time     `json:"last_updated"`
}

// Message types
const (
	MessageTypeTask        = "task"
	MessageTypeResult      = "result"
	MessageTypeStatus      = "status"
	MessageTypeEvent       = "event"
	MessageTypeHeartbeat   = "heartbeat"
	MessageTypeCommand     = "command"
	MessageTypeQuery       = "query"
	MessageTypeResponse    = "response"
	MessageTypeNotification = "notification"
)

// Event types
const (
	EventAgentStarted   = "agent_started"
	EventAgentStopped   = "agent_stopped"
	EventTaskCompleted  = "task_completed"
	EventTaskFailed     = "task_failed"
	EventStateChanged   = "state_changed"
	EventError          = "error"
	EventHealthCheck    = "health_check"
)

// Agent status constants
const (
	AgentStatusIdle        = "idle"
	AgentStatusBusy        = "busy"
	AgentStatusError       = "error"
	AgentStatusMaintenance = "maintenance"
	AgentStatusStopping    = "stopping"
)

// NewAgentCommunicationSystem creates a new communication system
func NewAgentCommunicationSystem(database *Database) *AgentCommunicationSystem {
	system := &AgentCommunicationSystem{
		database:      database,
		agents:        make(map[string]Agent),
		subscriptions: make(map[string][]string),
		middleware:    []Middleware{},
		metrics:       &CommunicationMetrics{
			QueueSizes:      make(map[string]int),
			ThroughputTrend: make([]float64, 0, 100),
		},
		messageQueue: &MessageQueue{
			queues:     make(map[string]*Queue),
			processing: make(map[string]*Message),
			workers:    10,
		},
		stateManager: &StateManager{
			states:    make(map[string]*AgentState),
			snapshots: make([]*StateSnapshot, 0, 100),
		},
		eventBus: &EventBus{
			subscribers: make(map[string][]EventHandler),
			events:      make([]*Event, 0, 1000),
		},
	}
	
	// Initialize message queues
	system.initializeQueues()
	
	// Setup default middleware
	system.setupDefaultMiddleware()
	
	// Create database tables
	if err := system.createCommunicationTables(); err != nil {
		log.Printf("Warning: Failed to create communication tables: %v", err)
	}
	
	// Start workers
	go system.startMessageWorkers()
	go system.startMetricsCollector()
	go system.startHealthChecker()
	
	log.Println("Agent Communication System initialized")
	
	return system
}

// RegisterAgent registers an agent with the communication system
func (acs *AgentCommunicationSystem) RegisterAgent(agent Agent) error {
	acs.mu.Lock()
	defer acs.mu.Unlock()
	
	agentID := agent.GetID()
	
	if _, exists := acs.agents[agentID]; exists {
		return fmt.Errorf("agent %s is already registered", agentID)
	}
	
	acs.agents[agentID] = agent
	
	// Initialize agent state
	state := agent.GetState()
	acs.stateManager.mu.Lock()
	acs.stateManager.states[agentID] = &state
	acs.stateManager.mu.Unlock()
	
	// Publish agent started event
	acs.publishEvent(&Event{
		ID:        generateEventID(),
		Type:      EventAgentStarted,
		Source:    "communication_system",
		Target:    "all",
		Data:      map[string]interface{}{"agent_id": agentID, "agent_type": agent.GetType()},
		Level:     "info",
		Timestamp: time.Now(),
	})
	
	log.Printf("Agent registered: %s (%s)", agentID, agent.GetType())
	return nil
}

// UnregisterAgent removes an agent from the communication system
func (acs *AgentCommunicationSystem) UnregisterAgent(agentID string) error {
	acs.mu.Lock()
	defer acs.mu.Unlock()
	
	if _, exists := acs.agents[agentID]; !exists {
		return fmt.Errorf("agent %s not found", agentID)
	}
	
	delete(acs.agents, agentID)
	delete(acs.subscriptions, agentID)
	
	// Remove agent state
	acs.stateManager.mu.Lock()
	delete(acs.stateManager.states, agentID)
	acs.stateManager.mu.Unlock()
	
	// Publish agent stopped event
	acs.publishEvent(&Event{
		ID:        generateEventID(),
		Type:      EventAgentStopped,
		Source:    "communication_system",
		Target:    "all",
		Data:      map[string]interface{}{"agent_id": agentID},
		Level:     "info",
		Timestamp: time.Now(),
	})
	
	log.Printf("Agent unregistered: %s", agentID)
	return nil
}

// SendMessage sends a message to an agent or broadcasts to all
func (acs *AgentCommunicationSystem) SendMessage(message *Message) error {
	// Generate message ID if not provided
	if message.ID == "" {
		message.ID = generateMessageID()
	}
	
	// Set timestamps
	message.CreatedAt = time.Now()
	if message.Metadata.TTL > 0 {
		expiresAt := message.CreatedAt.Add(message.Metadata.TTL)
		message.ExpiresAt = &expiresAt
	}
	
	// Apply middleware
	if err := acs.applyMiddleware(message); err != nil {
		return fmt.Errorf("middleware failed: %w", err)
	}
	
	// Validate message
	if err := acs.validateMessage(message); err != nil {
		return fmt.Errorf("message validation failed: %w", err)
	}
	
	// Route message
	if message.To == "broadcast" || message.To == "all" {
		return acs.broadcastMessage(message)
	} else {
		return acs.routeMessage(message)
	}
}

// Subscribe subscribes an agent to specific message types
func (acs *AgentCommunicationSystem) Subscribe(agentID string, messageTypes []string) error {
	acs.mu.Lock()
	defer acs.mu.Unlock()
	
	if _, exists := acs.agents[agentID]; !exists {
		return fmt.Errorf("agent %s not registered", agentID)
	}
	
	acs.subscriptions[agentID] = messageTypes
	
	log.Printf("Agent %s subscribed to message types: %v", agentID, messageTypes)
	return nil
}

// PublishEvent publishes an event to the event bus
func (acs *AgentCommunicationSystem) PublishEvent(event *Event) error {
	return acs.publishEvent(event)
}

// SubscribeToEvents subscribes to specific event types
func (acs *AgentCommunicationSystem) SubscribeToEvents(eventTypes []string, handler EventHandler) error {
	acs.eventBus.mu.Lock()
	defer acs.eventBus.mu.Unlock()
	
	for _, eventType := range eventTypes {
		acs.eventBus.subscribers[eventType] = append(acs.eventBus.subscribers[eventType], handler)
	}
	
	return nil
}

// UpdateAgentState updates an agent's state and broadcasts the change
func (acs *AgentCommunicationSystem) UpdateAgentState(agentID string, state AgentState) error {
	acs.stateManager.mu.Lock()
	defer acs.stateManager.mu.Unlock()
	
	oldState := acs.stateManager.states[agentID]
	state.UpdatedAt = time.Now()
	acs.stateManager.states[agentID] = &state
	
	// Publish state change event if significant change
	if oldState == nil || oldState.Status != state.Status || oldState.CurrentTask != state.CurrentTask {
		acs.publishEvent(&Event{
			ID:     generateEventID(),
			Type:   EventStateChanged,
			Source: agentID,
			Target: "all",
			Data: map[string]interface{}{
				"agent_id":  agentID,
				"old_state": oldState,
				"new_state": state,
			},
			Level:     "info",
			Timestamp: time.Now(),
		})
	}
	
	return nil
}

// GetAgentState retrieves the current state of an agent
func (acs *AgentCommunicationSystem) GetAgentState(agentID string) (*AgentState, error) {
	acs.stateManager.mu.RLock()
	defer acs.stateManager.mu.RUnlock()
	
	state, exists := acs.stateManager.states[agentID]
	if !exists {
		return nil, fmt.Errorf("agent %s not found", agentID)
	}
	
	return state, nil
}

// GetSystemState returns the state of all agents
func (acs *AgentCommunicationSystem) GetSystemState() map[string]*AgentState {
	acs.stateManager.mu.RLock()
	defer acs.stateManager.mu.RUnlock()
	
	// Create a copy of the states
	states := make(map[string]*AgentState)
	for agentID, state := range acs.stateManager.states {
		stateCopy := *state
		states[agentID] = &stateCopy
	}
	
	return states
}

// CreateStateSnapshot creates a snapshot of the current system state
func (acs *AgentCommunicationSystem) CreateStateSnapshot() *StateSnapshot {
	snapshot := &StateSnapshot{
		ID:        generateSnapshotID(),
		States:    acs.GetSystemState(),
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
	
	// Store snapshot
	acs.stateManager.mu.Lock()
	acs.stateManager.snapshots = append(acs.stateManager.snapshots, snapshot)
	
	// Keep only last 100 snapshots
	if len(acs.stateManager.snapshots) > 100 {
		acs.stateManager.snapshots = acs.stateManager.snapshots[1:]
	}
	acs.stateManager.mu.Unlock()
	
	return snapshot
}

// GetMetrics returns current communication metrics
func (acs *AgentCommunicationSystem) GetMetrics() *CommunicationMetrics {
	acs.updateMetrics()
	return acs.metrics
}

// AddMiddleware adds message processing middleware
func (acs *AgentCommunicationSystem) AddMiddleware(middleware Middleware) {
	acs.mu.Lock()
	defer acs.mu.Unlock()
	
	acs.middleware = append(acs.middleware, middleware)
}

// Internal methods

func (acs *AgentCommunicationSystem) initializeQueues() {
	priorities := []string{"critical", "high", "medium", "low"}
	
	for _, priority := range priorities {
		acs.messageQueue.queues[priority] = &Queue{
			Messages: make([]*Message, 0),
			Priority: priority,
			MaxSize:  1000,
		}
	}
}

func (acs *AgentCommunicationSystem) setupDefaultMiddleware() {
	// Logging middleware
	acs.middleware = append(acs.middleware, func(message *Message, next func(*Message) error) error {
		log.Printf("Processing message: %s -> %s (%s)", message.From, message.To, message.Type)
		return next(message)
	})
	
	// Metrics middleware
	acs.middleware = append(acs.middleware, func(message *Message, next func(*Message) error) error {
		start := time.Now()
		err := next(message)
		
		// Update metrics
		acs.mu.Lock()
		acs.metrics.TotalMessages++
		latency := time.Since(start)
		
		// Update average latency (simple moving average)
		if acs.metrics.AverageLatency == 0 {
			acs.metrics.AverageLatency = latency
		} else {
			acs.metrics.AverageLatency = (acs.metrics.AverageLatency + latency) / 2
		}
		
		if err != nil {
			acs.metrics.ErrorRate = (acs.metrics.ErrorRate*9 + 1) / 10 // Moving average
		} else {
			acs.metrics.ErrorRate = acs.metrics.ErrorRate * 0.9 // Decay error rate
		}
		acs.mu.Unlock()
		
		return err
	})
	
	// Persistence middleware
	acs.middleware = append(acs.middleware, func(message *Message, next func(*Message) error) error {
		if message.Metadata.Persistent {
			if err := acs.storeMessage(message); err != nil {
				log.Printf("Failed to store message: %v", err)
			}
		}
		return next(message)
	})
}

func (acs *AgentCommunicationSystem) applyMiddleware(message *Message) error {
	// Create middleware chain
	chain := func(msg *Message) error {
		return nil // End of chain
	}
	
	// Build chain in reverse order
	for i := len(acs.middleware) - 1; i >= 0; i-- {
		middleware := acs.middleware[i]
		next := chain
		chain = func(msg *Message) error {
			return middleware(msg, next)
		}
	}
	
	return chain(message)
}

func (acs *AgentCommunicationSystem) validateMessage(message *Message) error {
	if message.From == "" {
		return fmt.Errorf("message must have a sender")
	}
	
	if message.To == "" {
		return fmt.Errorf("message must have a recipient")
	}
	
	if message.Type == "" {
		return fmt.Errorf("message must have a type")
	}
	
	// Check if sender is registered
	acs.mu.RLock()
	_, exists := acs.agents[message.From]
	acs.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("sender agent %s is not registered", message.From)
	}
	
	// Check if recipient is registered (unless broadcast)
	if message.To != "broadcast" && message.To != "all" {
		acs.mu.RLock()
		_, exists := acs.agents[message.To]
		acs.mu.RUnlock()
		
		if !exists {
			return fmt.Errorf("recipient agent %s is not registered", message.To)
		}
	}
	
	return nil
}

func (acs *AgentCommunicationSystem) broadcastMessage(message *Message) error {
	acs.mu.RLock()
	agents := make([]Agent, 0, len(acs.agents))
	for _, agent := range acs.agents {
		agents = append(agents, agent)
	}
	acs.mu.RUnlock()
	
	var lastError error
	deliveredCount := 0
	
	for _, agent := range agents {
		// Skip sender
		if agent.GetID() == message.From {
			continue
		}
		
		// Check subscription
		if !acs.isSubscribed(agent.GetID(), message.Type) {
			continue
		}
		
		// Create copy for each recipient
		msgCopy := *message
		msgCopy.To = agent.GetID()
		
		if err := acs.queueMessage(&msgCopy); err != nil {
			lastError = err
			log.Printf("Failed to queue message to %s: %v", agent.GetID(), err)
		} else {
			deliveredCount++
		}
	}
	
	if deliveredCount == 0 && lastError != nil {
		return fmt.Errorf("failed to broadcast message: %w", lastError)
	}
	
	log.Printf("Broadcast message %s delivered to %d agents", message.ID, deliveredCount)
	return nil
}

func (acs *AgentCommunicationSystem) routeMessage(message *Message) error {
	return acs.queueMessage(message)
}

func (acs *AgentCommunicationSystem) queueMessage(message *Message) error {
	priority := message.Metadata.Priority
	if priority == "" {
		priority = "medium"
	}
	
	acs.messageQueue.mu.Lock()
	defer acs.messageQueue.mu.Unlock()
	
	queue, exists := acs.messageQueue.queues[priority]
	if !exists {
		return fmt.Errorf("unknown priority: %s", priority)
	}
	
	queue.mu.Lock()
	defer queue.mu.Unlock()
	
	// Check queue capacity
	if len(queue.Messages) >= queue.MaxSize {
		return fmt.Errorf("queue %s is full", priority)
	}
	
	// Check for urgent messages - add to front
	if message.Metadata.Urgent {
		queue.Messages = append([]*Message{message}, queue.Messages...)
	} else {
		queue.Messages = append(queue.Messages, message)
	}
	
	message.Status = "queued"
	
	return nil
}

func (acs *AgentCommunicationSystem) isSubscribed(agentID, messageType string) bool {
	acs.mu.RLock()
	defer acs.mu.RUnlock()
	
	subscriptions, exists := acs.subscriptions[agentID]
	if !exists {
		return true // By default, agents receive all messages
	}
	
	for _, subscribedType := range subscriptions {
		if subscribedType == messageType || subscribedType == "*" {
			return true
		}
	}
	
	return false
}

func (acs *AgentCommunicationSystem) publishEvent(event *Event) error {
	if event.ID == "" {
		event.ID = generateEventID()
	}
	
	event.Timestamp = time.Now()
	
	// Store event
	acs.eventBus.mu.Lock()
	acs.eventBus.events = append(acs.eventBus.events, event)
	
	// Keep only last 1000 events
	if len(acs.eventBus.events) > 1000 {
		acs.eventBus.events = acs.eventBus.events[1:]
	}
	
	// Get subscribers
	subscribers := make([]EventHandler, 0)
	if handlers, exists := acs.eventBus.subscribers[event.Type]; exists {
		subscribers = append(subscribers, handlers...)
	}
	if handlers, exists := acs.eventBus.subscribers["*"]; exists {
		subscribers = append(subscribers, handlers...)
	}
	acs.eventBus.mu.Unlock()
	
	// Notify subscribers
	for _, handler := range subscribers {
		go func(h EventHandler, e *Event) {
			if err := h(e); err != nil {
				log.Printf("Event handler failed: %v", err)
			}
		}(handler, event)
	}
	
	return nil
}

func (acs *AgentCommunicationSystem) startMessageWorkers() {
	for i := 0; i < acs.messageQueue.workers; i++ {
		go acs.messageWorker(i)
	}
}

func (acs *AgentCommunicationSystem) messageWorker(workerID int) {
	log.Printf("Message worker %d started", workerID)
	
	for {
		message := acs.getNextMessage()
		if message == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		
		// Check expiration
		if message.ExpiresAt != nil && time.Now().After(*message.ExpiresAt) {
			message.Status = "expired"
			log.Printf("Message %s expired", message.ID)
			continue
		}
		
		// Mark as processing
		acs.messageQueue.mu.Lock()
		acs.messageQueue.processing[message.ID] = message
		acs.messageQueue.mu.Unlock()
		
		// Deliver message
		if err := acs.deliverMessage(message); err != nil {
			log.Printf("Worker %d failed to deliver message %s: %v", workerID, message.ID, err)
			
			// Handle retry
			if message.RetryCount < message.MaxRetries {
				message.RetryCount++
				message.Status = "retry"
				acs.queueMessage(message)
			} else {
				message.Status = "failed"
			}
		} else {
			message.Status = "delivered"
			now := time.Now()
			message.DeliveredAt = &now
		}
		
		// Remove from processing
		acs.messageQueue.mu.Lock()
		delete(acs.messageQueue.processing, message.ID)
		acs.messageQueue.mu.Unlock()
	}
}

func (acs *AgentCommunicationSystem) getNextMessage() *Message {
	acs.messageQueue.mu.RLock()
	defer acs.messageQueue.mu.RUnlock()
	
	// Check queues in priority order
	priorities := []string{"critical", "high", "medium", "low"}
	
	for _, priority := range priorities {
		queue := acs.messageQueue.queues[priority]
		
		queue.mu.Lock()
		if len(queue.Messages) > 0 {
			message := queue.Messages[0]
			queue.Messages = queue.Messages[1:]
			queue.mu.Unlock()
			return message
		}
		queue.mu.Unlock()
	}
	
	return nil
}

func (acs *AgentCommunicationSystem) deliverMessage(message *Message) error {
	acs.mu.RLock()
	agent, exists := acs.agents[message.To]
	acs.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("recipient agent %s not found", message.To)
	}
	
	return agent.HandleMessage(message)
}

func (acs *AgentCommunicationSystem) startMetricsCollector() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	
	for range ticker.C {
		acs.updateMetrics()
	}
}

func (acs *AgentCommunicationSystem) updateMetrics() {
	acs.mu.Lock()
	defer acs.mu.Unlock()
	
	// Update queue sizes
	for priority, queue := range acs.messageQueue.queues {
		queue.mu.RLock()
		acs.metrics.QueueSizes[priority] = len(queue.Messages)
		queue.mu.RUnlock()
	}
	
	// Update agent counts
	acs.metrics.ActiveAgents = len(acs.agents)
	
	healthyCount := 0
	for _, agent := range acs.agents {
		if agent.IsHealthy() {
			healthyCount++
		}
	}
	acs.metrics.HealthyAgents = healthyCount
	
	// Calculate messages per second
	now := time.Now()
	if !acs.metrics.LastUpdated.IsZero() {
		timeDiff := now.Sub(acs.metrics.LastUpdated).Seconds()
		if timeDiff > 0 {
			acs.metrics.MessagesPerSecond = float64(acs.metrics.TotalMessages) / timeDiff
		}
	}
	
	// Update throughput trend
	acs.metrics.ThroughputTrend = append(acs.metrics.ThroughputTrend, acs.metrics.MessagesPerSecond)
	if len(acs.metrics.ThroughputTrend) > 100 {
		acs.metrics.ThroughputTrend = acs.metrics.ThroughputTrend[1:]
	}
	
	acs.metrics.LastUpdated = now
}

func (acs *AgentCommunicationSystem) startHealthChecker() {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	
	for range ticker.C {
		acs.performHealthCheck()
	}
}

func (acs *AgentCommunicationSystem) performHealthCheck() {
	acs.mu.RLock()
	agents := make([]Agent, 0, len(acs.agents))
	for _, agent := range acs.agents {
		agents = append(agents, agent)
	}
	acs.mu.RUnlock()
	
	for _, agent := range agents {
		if !agent.IsHealthy() {
			acs.publishEvent(&Event{
				ID:     generateEventID(),
				Type:   EventError,
				Source: "communication_system",
				Target: "all",
				Data:   map[string]interface{}{"agent_id": agent.GetID(), "error": "health check failed"},
				Level:  "warning",
				Timestamp: time.Now(),
			})
		}
	}
}

// Database operations
func (acs *AgentCommunicationSystem) createCommunicationTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			from_agent TEXT NOT NULL,
			to_agent TEXT NOT NULL,
			subject TEXT,
			payload TEXT NOT NULL,
			metadata TEXT NOT NULL,
			status TEXT NOT NULL,
			retry_count INTEGER NOT NULL,
			created_at INTEGER NOT NULL,
			delivered_at INTEGER,
			expires_at INTEGER
		)`,
		`CREATE TABLE IF NOT EXISTS agent_states (
			agent_id TEXT PRIMARY KEY,
			status TEXT NOT NULL,
			current_task TEXT,
			progress REAL NOT NULL,
			state_data TEXT NOT NULL,
			last_activity INTEGER NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			source TEXT NOT NULL,
			target TEXT NOT NULL,
			data TEXT NOT NULL,
			level TEXT NOT NULL,
			timestamp INTEGER NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_status ON messages(status)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_type ON messages(type)`,
		`CREATE INDEX IF NOT EXISTS idx_agent_states_status ON agent_states(status)`,
		`CREATE INDEX IF NOT EXISTS idx_events_type ON events(type)`,
		`CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp)`,
	}
	
	for _, query := range queries {
		if _, err := acs.database.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create communication table: %w", err)
		}
	}
	
	return nil
}

func (acs *AgentCommunicationSystem) storeMessage(message *Message) error {
	payloadJSON, err := json.Marshal(message.Payload)
	if err != nil {
		return err
	}
	
	metadataJSON, err := json.Marshal(message.Metadata)
	if err != nil {
		return err
	}
	
	var deliveredAt, expiresAt *int64
	if message.DeliveredAt != nil {
		t := message.DeliveredAt.Unix()
		deliveredAt = &t
	}
	if message.ExpiresAt != nil {
		t := message.ExpiresAt.Unix()
		expiresAt = &t
	}
	
	query := `
		INSERT OR REPLACE INTO messages
		(id, type, from_agent, to_agent, subject, payload, metadata, status, 
		 retry_count, created_at, delivered_at, expires_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err = acs.database.db.Exec(query,
		message.ID, message.Type, message.From, message.To, message.Subject,
		string(payloadJSON), string(metadataJSON), message.Status,
		message.RetryCount, message.CreatedAt.Unix(), deliveredAt, expiresAt)
	
	return err
}

// Public API methods

// GetMessageHistory retrieves message history
func (acs *AgentCommunicationSystem) GetMessageHistory(agentID string, limit int) ([]*Message, error) {
	if limit <= 0 {
		limit = 100
	}
	
	query := `
		SELECT id, type, from_agent, to_agent, subject, payload, metadata, status,
		       retry_count, created_at, delivered_at, expires_at
		FROM messages 
		WHERE from_agent = ? OR to_agent = ?
		ORDER BY created_at DESC 
		LIMIT ?
	`
	
	rows, err := acs.database.db.Query(query, agentID, agentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var messages []*Message
	for rows.Next() {
		var message Message
		var payloadJSON, metadataJSON string
		var deliveredAtUnix, expiresAtUnix *int64
		
		err := rows.Scan(&message.ID, &message.Type, &message.From, &message.To,
			&message.Subject, &payloadJSON, &metadataJSON, &message.Status,
			&message.RetryCount, &message.CreatedAt, &deliveredAtUnix, &expiresAtUnix)
		if err != nil {
			continue
		}
		
		json.Unmarshal([]byte(payloadJSON), &message.Payload)
		json.Unmarshal([]byte(metadataJSON), &message.Metadata)
		
		if deliveredAtUnix != nil {
			t := time.Unix(*deliveredAtUnix, 0)
			message.DeliveredAt = &t
		}
		if expiresAtUnix != nil {
			t := time.Unix(*expiresAtUnix, 0)
			message.ExpiresAt = &t
		}
		
		messages = append(messages, &message)
	}
	
	return messages, nil
}

// GetEventHistory retrieves event history
func (acs *AgentCommunicationSystem) GetEventHistory(limit int) ([]*Event, error) {
	acs.eventBus.mu.RLock()
	defer acs.eventBus.mu.RUnlock()
	
	if limit <= 0 || limit > len(acs.eventBus.events) {
		limit = len(acs.eventBus.events)
	}
	
	start := len(acs.eventBus.events) - limit
	events := make([]*Event, limit)
	copy(events, acs.eventBus.events[start:])
	
	return events, nil
}

// GetStateHistory retrieves state snapshots
func (acs *AgentCommunicationSystem) GetStateHistory(limit int) ([]*StateSnapshot, error) {
	acs.stateManager.mu.RLock()
	defer acs.stateManager.mu.RUnlock()
	
	if limit <= 0 || limit > len(acs.stateManager.snapshots) {
		limit = len(acs.stateManager.snapshots)
	}
	
	start := len(acs.stateManager.snapshots) - limit
	snapshots := make([]*StateSnapshot, limit)
	copy(snapshots, acs.stateManager.snapshots[start:])
	
	return snapshots, nil
}

// GetRegisteredAgents returns all registered agents
func (acs *AgentCommunicationSystem) GetRegisteredAgents() map[string]Agent {
	acs.mu.RLock()
	defer acs.mu.RUnlock()
	
	agents := make(map[string]Agent)
	for id, agent := range acs.agents {
		agents[id] = agent
	}
	
	return agents
}

// Shutdown gracefully shuts down the communication system
func (acs *AgentCommunicationSystem) Shutdown() error {
	log.Println("Shutting down Agent Communication System...")
	
	// Create final state snapshot
	acs.CreateStateSnapshot()
	
	// Notify all agents of shutdown
	shutdownMessage := &Message{
		ID:      generateMessageID(),
		Type:    MessageTypeCommand,
		From:    "communication_system",
		To:      "broadcast",
		Subject: "shutdown",
		Payload: map[string]interface{}{"command": "shutdown"},
		Metadata: MessageMetadata{
			Priority: "critical",
			Urgent:   true,
		},
	}
	
	acs.broadcastMessage(shutdownMessage)
	
	// Wait a moment for messages to be processed
	time.Sleep(time.Second * 2)
	
	log.Println("Agent Communication System shut down complete")
	return nil
}

// Helper functions

func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

func generateSnapshotID() string {
	return fmt.Sprintf("snap_%d", time.Now().UnixNano())
}