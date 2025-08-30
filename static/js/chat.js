class ChatClient {
    constructor() {
        this.ws = null;
        this.isConnected = false;
        this.userID = null;
        this.username = null;
        
        this.initializeEventListeners();
    }

    initializeEventListeners() {
        // 回车发送消息
        document.getElementById('message-input').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.sendMessage();
            }
        });
    }

    connect() {
        const userID = document.getElementById('user-id').value;
        const username = document.getElementById('username').value;

        if (!userID || !username) {
            alert('请输入用户ID和用户名');
            return;
        }

        this.userID = parseInt(userID);
        this.username = username;

        const wsUrl = `ws://localhost:8081/ws?user_id=${userID}&username=${encodeURIComponent(username)}`;
        
        try {
            this.ws = new WebSocket(wsUrl);
            
            this.ws.onopen = (event) => {
                this.onConnect(event);
            };
            
            this.ws.onmessage = (event) => {
                this.onMessage(event);
            };
            
            this.ws.onclose = (event) => {
                this.onDisconnect(event);
            };
            
            this.ws.onerror = (error) => {
                this.onError(error);
            };
            
        } catch (error) {
            console.error('WebSocket连接失败:', error);
            this.updateStatus('连接失败', 'disconnected');
        }
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
    }

    onConnect(event) {
        console.log('WebSocket连接已建立');
        this.isConnected = true;
        this.updateStatus('已连接', 'connected');
        this.updateUI();
        this.addSystemMessage('已连接到聊天服务器');
        
        // 发送心跳
        this.startHeartbeat();
    }

    onMessage(event) {
        try {
            const message = JSON.parse(event.data);
            this.handleMessage(message);
        } catch (error) {
            console.error('消息解析失败:', error);
        }
    }

    onDisconnect(event) {
        console.log('WebSocket连接已断开');
        this.isConnected = false;
        this.updateStatus('已断开', 'disconnected');
        this.updateUI();
        this.addSystemMessage('与聊天服务器的连接已断开');
        
        // 停止心跳
        this.stopHeartbeat();
    }

    onError(error) {
        console.error('WebSocket错误:', error);
        this.updateStatus('连接错误', 'disconnected');
    }

    handleMessage(message) {
        console.log('收到消息:', message);

        switch (message.type) {
            case 'private':
                this.handlePrivateMessage(message);
                break;
            case 'user_list':
                this.handleUserList(message);
                break;
            case 'join':
                this.handleUserJoin(message);
                break;
            case 'leave':
                this.handleUserLeave(message);
                break;
            case 'error':
                this.handleError(message);
                break;
            case 'heartbeat':
                // 心跳响应，不需要处理
                break;
            case 'read':
                this.handleReadReceipt(message);
                break;
            default:
                console.log('未知消息类型:', message.type);
        }
    }

    handlePrivateMessage(message) {
        this.addMessage({
            from: message.from_user_id,
            content: message.content,
            timestamp: message.timestamp,
            type: 'received'
        });
    }

    handleUserList(message) {
        this.updateOnlineUsers(message.data.users);
    }

    handleUserJoin(message) {
        const user = message.data;
        this.addSystemMessage(`${user.username} 加入了聊天`);
        this.updateOnlineUserCount();
    }

    handleUserLeave(message) {
        const user = message.data;
        this.addSystemMessage(`${user.username} 离开了聊天`);
        this.updateOnlineUserCount();
    }

    handleError(message) {
        const errorData = message.data;
        this.addSystemMessage(`错误: ${errorData.message}`, 'error');
    }

    handleReadReceipt(message) {
        if (message.data && message.data.status === 'delivered') {
            console.log('消息已送达');
        }
    }

    sendMessage() {
        if (!this.isConnected) {
            alert('请先连接到服务器');
            return;
        }

        const messageInput = document.getElementById('message-input');
        const targetUserInput = document.getElementById('target-user');
        
        const content = messageInput.value.trim();
        const targetUserID = parseInt(targetUserInput.value);

        if (!content) {
            alert('请输入消息内容');
            return;
        }

        if (!targetUserID) {
            alert('请输入目标用户ID');
            return;
        }

        const message = {
            type: 'private',
            to_user_id: targetUserID,
            content: content,
            timestamp: new Date().toISOString()
        };

        try {
            this.ws.send(JSON.stringify(message));
            
            // 添加到消息列表
            this.addMessage({
                from: this.userID,
                to: targetUserID,
                content: content,
                timestamp: new Date().toISOString(),
                type: 'sent'
            });
            
            // 清空输入框
            messageInput.value = '';
            
        } catch (error) {
            console.error('发送消息失败:', error);
            alert('发送消息失败');
        }
    }

    addMessage(messageData) {
        const messagesContainer = document.getElementById('messages');
        const messageElement = document.createElement('div');
        
        const time = new Date(messageData.timestamp).toLocaleTimeString();
        const messageClass = messageData.type === 'sent' ? 'sent' : 'received';
        
        messageElement.className = `message ${messageClass}`;
        messageElement.innerHTML = `
            <div class="message-header">
                <span>用户 ${messageData.from || '我'}</span>
                <span>${time}</span>
            </div>
            <div class="message-content">${this.escapeHtml(messageData.content)}</div>
        `;
        
        messagesContainer.appendChild(messageElement);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    addSystemMessage(content, type = 'system') {
        const messagesContainer = document.getElementById('messages');
        const messageElement = document.createElement('div');
        
        messageElement.className = `message ${type}`;
        messageElement.innerHTML = `
            <div class="message-content">${this.escapeHtml(content)}</div>
        `;
        
        messagesContainer.appendChild(messageElement);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    updateOnlineUsers(users) {
        const usersContainer = document.getElementById('online-users');
        usersContainer.innerHTML = '';
        
        users.forEach(user => {
            const userElement = document.createElement('div');
            userElement.className = 'user-item';
            userElement.innerHTML = `
                <div>${this.escapeHtml(user.username)}</div>
                <div style="font-size: 12px; color: #666;">ID: ${user.user_id}</div>
            `;
            
            // 点击用户名自动填入目标用户ID
            userElement.addEventListener('click', () => {
                document.getElementById('target-user').value = user.user_id;
            });
            
            usersContainer.appendChild(userElement);
        });
        
        this.updateOnlineUserCount(users.length);
    }

    updateOnlineUserCount(count) {
        if (count === undefined) {
            // 通过API获取在线用户数
            fetch('/api/v1/ws/online-users')
                .then(response => response.json())
                .then(data => {
                    console.log('在线用户API响应:', data);
                    if (data.code === 200) {
                        document.getElementById('online-count').textContent = `在线用户: ${data.data.count}`;
                    }
                })
                .catch(error => console.error('获取在线用户数失败:', error));
        } else {
            document.getElementById('online-count').textContent = `在线用户: ${count}`;
        }
    }

    updateStatus(status, className) {
        const statusElement = document.getElementById('status');
        statusElement.textContent = status;
        statusElement.className = `status ${className}`;
    }

    updateUI() {
        const connectBtn = document.getElementById('connect-btn');
        const disconnectBtn = document.getElementById('disconnect-btn');
        const messageInput = document.getElementById('message-input');
        const sendBtn = document.getElementById('send-btn');
        
        connectBtn.disabled = this.isConnected;
        disconnectBtn.disabled = !this.isConnected;
        messageInput.disabled = !this.isConnected;
        sendBtn.disabled = !this.isConnected;
    }

    startHeartbeat() {
        this.heartbeatInterval = setInterval(() => {
            if (this.isConnected && this.ws.readyState === WebSocket.OPEN) {
                const heartbeat = {
                    type: 'heartbeat',
                    content: 'ping',
                    timestamp: new Date().toISOString()
                };
                this.ws.send(JSON.stringify(heartbeat));
            }
        }, 30000); // 每30秒发送一次心跳
    }

    stopHeartbeat() {
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
            this.heartbeatInterval = null;
        }
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// 全局变量和函数
let chatClient;

function init() {
    chatClient = new ChatClient();
}

function connect() {
    chatClient.connect();
}

function disconnect() {
    chatClient.disconnect();
}

function sendMessage() {
    chatClient.sendMessage();
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', init);
