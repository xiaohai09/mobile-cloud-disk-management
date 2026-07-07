# User Guide

[![Vue 3](https://img.shields.io/badge/Vue-3.4-green.svg)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.3-blue.svg)](https://www.typescriptlang.org/)
[![Element Plus](https://img.shields.io/badge/Element_Plus-2.4-blue.svg)](https://element-plus.org/)

> This guide is for end users, covering role descriptions, feature module operations, common troubleshooting, and best practices.

---

## 📑 Table of Contents

- [User Roles and Permissions](#user-roles-and-permissions)
- [Getting Started](#getting-started)
- [Feature Modules](#feature-modules)
- [Operation Workflow Examples](#operation-workflow-examples)
- [Common Troubleshooting](#common-troubleshooting)
- [Best Practices](#best-practices)
- [Glossary](#glossary)
- [Technical Support](#technical-support)

---

## User Roles and Permissions

The system uses Role-Based Access Control (RBAC) with three user roles:

| Role | Permissions | Use Case |
|------|-------------|----------|
| **Regular User** | View own accounts, tasks, redemption records | Daily operations staff |
| **Administrator** | Manage all users, tasks, redemption config, announcements | Team leads |
| **Super Administrator** | System configuration, audit logs, all admin permissions | IT administrators |

### Permission Matrix

| Feature Module | Regular User | Administrator | Super Administrator |
|---------------|--------------|---------------|---------------------|
| Dashboard | View | View | View |
| Account Management | View/edit own accounts | View/edit all accounts | View/edit all accounts |
| Task Management | View/manage own tasks | View/manage all tasks | View/manage all tasks |
| Redemption Center | View/create redemptions | View/manage all redemptions | View/manage all redemptions |
| Data Export | Export own data | Export all data | Export all data |
| Webhook Management | View own webhooks | View/manage all webhooks | View/manage all webhooks |
| System Settings | Edit personal info | View system settings | Full system configuration |
| Audit Logs | No access | View operation logs | View/export audit logs |

---

## Getting Started

### First Login

1. Visit the system homepage: `http://localhost`
2. Log in with the default admin account:
   - Username: `admin`
   - Password: Must be changed on first login
3. After login, enter the dashboard to view system overview

### Change Password

1. Click avatar in the top right corner → Personal Settings
2. Select "Change Password"
3. Enter current password and new password
4. Click "Confirm"

> **Security Note**: Use strong passwords with at least 8 characters, including uppercase, lowercase, numbers, and special characters.

### Interface Overview

```
┌─────────────────────────────────────────────────────────────┐
│  [Logo] Mobile Cloud Disk Management    [Avatar] [Notify] [Settings] │
├──────────┬──────────────────────────────────────────────────┤
│ Dashboard │                                                  │
│ Accounts  │            Main Content Area                     │
│ Tasks     │                                                  │
│ Redemption│                                                  │
│ Export    │                                                  │
│ Webhook   │                                                  │
│ Settings  │                                                  │
└──────────┴──────────────────────────────────────────────────┘
```

---

## Feature Modules

### 📊 Dashboard

System homepage providing key metrics overview.

**Features**：

| Metric | Description |
|--------|-------------|
| System Health | API, Worker, database, Redis connection status |
| Total Accounts | Number of added cloud disk accounts |
| Today's Sign-ins | Today's automatic sign-in success/failure count |
| Task Queue | Pending, running, completed task counts |
| Redemption Success Rate | 7-day redemption success rate trend |
| Announcements | Latest system announcements and notifications |

**Operations**：
- Click metric cards for details
- Click chart areas to drill down data
- Refresh button to get latest data

---

### 👤 Account Management

Manage multiple cloud disk service accounts.

**Features**：

| Operation | Description |
|-----------|-------------|
| Add Account | Enter cloud disk account information (username, password, platform) |
| Edit Account | Modify account information, status, remarks |
| Delete Account | Delete unused accounts |
| Batch Import | Support CSV batch account import |
| Status Monitoring | View account online status, last sign-in time |

**Workflow**：

1. Click "Account Management" in the left menu
2. Click "Add Account" button
3. Fill in account information:
   - Platform: Select cloud disk platform type
   - Username: Cloud disk account username
   - Password: Cloud disk account password
   - Remark: Optional remark information
4. Click "Save"

**Batch Import CSV Example**：

```csv
platform,username,password,remark
caiyun,user001,pass001,Account group A
caiyun,user002,pass002,Account group A
caiyun,user003,pass003,Account group B
```

---

### ⏰ Task Center

Manage scheduled tasks and automated operations.

**Features**：

| Operation | Description |
|-----------|-------------|
| Create Task | Set task name, type, execution time, associated accounts |
| View Tasks | View task list, status, execution history |
| Start/Stop | Enable/disable task execution |
| Delete Task | Delete unneeded tasks |

**Task Types**：

| Type | Description | Cron Example |
|------|-------------|--------------|
| Sign-in Task | Automatic cloud disk sign-in | `0 8 * * *` (daily at 8 AM) |
| Redemption Task | Automatic product redemption | `0 10 * * *` (daily at 10 AM) |
| Check Task | Check account status | `0 */6 * * *` (every 6 hours) |

**Workflow**：

1. Click "Task Center" in the left menu
2. Click "Create Task" button
3. Fill in task information:
   - Task Name: Custom task name
   - Task Type: Select sign-in/redemption/check
   - Execution Time: Cron expression
   - Associated Accounts: Select accounts to associate
4. Click "Save"

---

### 🎁 Redemption Center

Manage redemption codes and product exchanges.

**Features**：

| Operation | Description |
|-----------|-------------|
| Browse Products | View available products for redemption |
| Create Redemption | Select account and product to create redemption task |
| View Records | View redemption history |
| Export Records | Export redemption records as CSV/JSON |

**Workflow**：

1. Click "Redemption Center" in the left menu
2. Browse available product list
3. Select product to redeem
4. Select account to execute redemption
5. Click "Redeem Now"

---

### 📤 Export Center

Export system data in CSV/JSON format.

**Features**：

| Operation | Description |
|-----------|-------------|
| Create Export | Select data type, time range, account filters |
| Download Files | Download completed export files |
| View History | View export task history |
| Delete Records | Delete expired export records |

**Export Types**：

| Type | Description | Field Examples |
|------|-------------|----------------|
| Redemption Records | Export redemption history | Redemption time, account, product, status |
| Task Logs | Export task execution logs | Task name, execution time, result |
| Account Information | Export basic account info | Account, platform, status, last sign-in time |

**Workflow**：

1. Click "Export Center" in the left menu
2. Click "Create Export" button
3. Select export type
4. Set time range (optional)
5. Select account filters (optional)
6. Choose export format: CSV or JSON
7. Click "Start Export"
8. Wait for task completion and download file

---

### 🔔 Webhook Management

Configure event notifications and external system integration.

**Features**：

| Operation | Description |
|-----------|-------------|
| Create Endpoint | Create Webhook receiver endpoint |
| Subscribe Events | Select event types to subscribe to |
| View Logs | View Webhook delivery history |
| Test Delivery | Manually trigger test delivery |

**Supported Event Types**：

| Event Type | Trigger |
|-----------|---------|
| `exchange.completed` | Exchange task completed |
| `exchange.failed` | Exchange task failed |
| `exchange.created` | Exchange task created |
| `task.completed` | Task execution completed |
| `task.failed` | Task execution failed |
| `account.created` | Account created |
| `account.updated` | Account updated |
| `account.deleted` | Account deleted |

**Workflow**：

1. Click "Webhook Management" in the left menu
2. Click "Create Endpoint" button
3. Fill in endpoint information:
   - Name: Custom endpoint name
   - URL: Webhook receiver URL (HTTPS recommended)
   - Events: Select event types to subscribe to
4. Copy generated Secret (displayed only once)
5. Click "Save"
6. Configure signature verification in target system

---

### ⚙️ System Settings

Manage personal information and system preferences.

**Features**：

| Setting | Description |
|---------|-------------|
| Personal Information | Modify nickname, avatar, contact info |
| Password Change | Change login password |
| Theme Settings | Light/Dark/Auto theme toggle |
| Language Settings | Chinese/English |
| Notification Settings | System notifications, email preferences |

---

## Operation Workflow Examples

### Example 1: Batch Import Accounts and Create Sign-in Task

1. **Prepare account data**

   Create `accounts.csv` file：
   ```csv
   platform,username,password,remark
   caiyun,user001,pass001,Group A
   caiyun,user002,pass002,Group A
   caiyun,user003,pass003,Group B
   ```

2. **Batch import accounts**
   - Go to "Account Management"
   - Click "Batch Import"
   - Upload CSV file
   - Confirm import

3. **Create sign-in task**
   - Go to "Task Center"
   - Click "Create Task"
   - Task Name: Daily Sign-in
   - Task Type: Sign-in Task
   - Cron Expression: `0 8 * * *` (daily at 8 AM)
   - Associated Accounts: Select all accounts in Group A
   - Click "Save"

4. **Enable task**
   - Find the newly created task in task list
   - Click "Enable" button

---

### Example 2: Configure Webhook Notifications

1. **Create Webhook endpoint**
   - Go to "Webhook Management"
   - Click "Create Endpoint"
   - Name: Redemption Completion Notification
   - URL: `https://your-server.com/webhook`
   - Events: Select `exchange.completed`, `exchange.failed`
   - Click "Save"
   - **Copy Secret** (displayed only once)

2. **Configure receiver**

   ```python
   # Flask receiver example
   from flask import Flask, request
   import hmac
   import hashlib

   app = Flask(__name__)
   SECRET = "whsec_xxxxxxxxxxxxxxxxxxxx"

   @app.route('/webhook', methods=['POST'])
   def webhook():
       signature = request.headers.get('X-Caiyun-Signature')
       timestamp = request.headers.get('X-Caiyun-Timestamp')
       body = request.get_data(as_text=True)

       # Verify signature
       expected = hmac.new(
           SECRET.encode(),
           f"{timestamp}\n{body}".encode(),
           hashlib.sha256
       ).hexdigest()

       if not hmac.compare_digest(expected, signature):
           return 'Invalid signature', 403

       # Process event
       event = request.json
       print(f"Received event: {event}")
       return 'OK', 200
   ```

3. **Test delivery**
   - Click "Test" in Webhook list
   - Check delivery log for success confirmation

---

## Common Troubleshooting

### 1. Account Addition Failure

**Possible causes**:
- Incorrect account password
- Network connection timeout
- Cloud disk service rate limiting

**Solutions**:
- Check if account password is correct
- Retry later
- Check task logs for detailed error information

---

### 2. Redemption Task Not Executed

**Possible causes**:
- Task not enabled
- Queue busy at execution time
- Account Cookie expired

**Solutions**:
- Check if task status is "Enabled"
- Check queue monitoring page
- Update account Cookie

---

### 3. Export File Download Failure

**Possible causes**:
- Export task not completed
- File expired (7-day auto cleanup)
- Insufficient disk space

**Solutions**:
- Wait for task completion before downloading
- Recreate export task
- Contact administrator to clean up disk

---

### 4. Webhook Reception Failure

**Possible causes**:
- Target URL unreachable
- HMAC signature verification failed
- Non-2xx response status code

**Solutions**:
- Check if target service is accessible
- Verify Webhook Secret configuration
- Check detailed error information in delivery logs

---

### 5. Task Execution Timeout

**Possible causes**:
- High network latency
- Slow target service response
- Too many concurrent tasks

**Solutions**:
- Adjust task execution time to avoid peak hours
- Reduce concurrent task count
- Contact administrator to check system resources

---

### 6. Login Failure

**Possible causes**:
- Incorrect username or password
- Account locked (multiple failed login attempts)
- System under maintenance

**Solutions**:
- Confirm correct username and password
- Wait for account unlock (usually 15 minutes)
- Contact administrator to confirm system status

---

## Best Practices

### Account Management

- Add clear remarks for accounts with different purposes
- Regularly check account status, handle abnormal accounts promptly
- Do not reuse passwords, use password manager to generate strong passwords

### Task Management

- Set reasonable task execution times to avoid peak hours
- Enable failure retry for important tasks
- Regularly check task execution logs to discover anomalies early

### Data Security

- Change login passwords regularly
- Do not save login status on public computers
- Export important data and store securely in time

### Webhook Usage

- Use HTTPS endpoints
- Keep Webhook Secret safe
- Implement idempotency handling to prevent duplicate deliveries

---

## Glossary

| Term | Description |
|------|-------------|
| Account | Cloud disk service login credentials, including username, password, platform info |
| Task | Scheduled automated operations, such as sign-in, redemption |
| Redemption | Use redemption codes to redeem cloud disk service products |
| Webhook | Event notification mechanism, system actively pushes events to external URLs |
| Cron | Scheduled task expression, defines task execution time |
| Token | Identity authentication token for API request authentication |
| JWT | JSON Web Token, an open standard token format |

---

## Technical Support

### Self-Service Support

- **System Logs**: `docker compose logs -f`
- **Monitoring Dashboard**: Access http://localhost:3000
- **API Documentation**: See [API Documentation](API.md)
- **Deployment Guide**: See [Deployment Guide](DEPLOYMENT.md)

### Contact Support

- **Issue Feedback**: [GitHub Issues](https://github.com/xiaohai09/mobile-cloud-disk-management/issues)
- **Security Vulnerabilities**: [Security Advisories](https://github.com/xiaohai09/mobile-cloud-disk-management/security/advisories)
- **Discussions**: [GitHub Discussions](https://github.com/xiaohai09/mobile-cloud-disk-management/discussions)

### Feedback Requirements

When submitting an Issue, please include：
1. Problem description
2. Reproduction steps
3. Expected behavior
4. Actual behavior
5. Environment information (browser, OS, etc.)
6. Screenshots or logs (if applicable)

---

*Last updated: 2026-07-07*
