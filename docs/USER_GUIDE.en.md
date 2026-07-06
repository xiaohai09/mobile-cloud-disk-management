# User Guide

## User Roles

The system supports three user roles:

| Role | Permissions |
|------|-------------|
| **Regular User** | View own accounts, tasks, exchange records |
| **Admin** | Manage all users, tasks, exchange configs, announcements |
| **Super Admin** | System configuration, audit logs, all admin permissions |

## First Login

1. Access the system homepage
2. Log in with default admin credentials:
   - Username: `admin`
   - Password: Must be changed on first login
3. After login, enter the dashboard to view system overview

## Feature Modules

### Dashboard

- System health status overview
- Task queue monitoring
- Exchange success rate statistics
- Announcement notification center

### Account Management

- Add/edit/delete cloud disk accounts
- Bulk import accounts
- Account status monitoring (online/offline/abnormal)
- Automatic sign-in status view

### Task Center

- Create exchange tasks
- View task execution logs
- Task start/stop control
- Task execution result statistics

### Exchange Center

- Product browsing and filtering
- Create exchange tasks
- Exchange record queries
- Exchange result export

### Export Center

- Export exchange records, task logs, account information
- Export formats: CSV, JSON
- Filter by time range and accounts
- Export history management

### Webhook Management

- Create webhook endpoints
- Subscribe to event types
- View delivery logs
- HMAC signature verification configuration

### System Settings

- Modify personal information
- Change avatar
- Theme switching (light/dark/auto)
- Password modification

## Common Issues

### 1. Account Addition Failed

**Possible causes**:
- Incorrect account password
- Network connection timeout
- Cloud disk server rate limiting

**Solutions**:
- Check if account credentials are correct
- Retry later
- Check task logs for detailed error information

### 2. Exchange Task Not Executing

**Possible causes**:
- Task not enabled
- Queue busy at execution time
- Account cookie expired

**Solutions**:
- Check if task status is "Enabled"
- Check queue monitoring page
- Update account cookie

### 3. Export File Download Failed

**Possible causes**:
- Export task not completed
- File expired (auto-cleaned after 7 days)
- Insufficient disk space

**Solutions**:
- Wait for task completion before downloading
- Recreate export task
- Contact admin to clean up disk

### 4. Webhook Reception Failed

**Possible causes**:
- Target URL unreachable
- HMAC signature verification failed
- Response status code not 2xx

**Solutions**:
- Check if target service is accessible
- Verify webhook secret configuration
- Check detailed error information in delivery logs

## Technical Support

For other issues:

1. Check system logs: `docker compose logs -f`
2. Check Grafana monitoring dashboard
3. Contact system administrator
