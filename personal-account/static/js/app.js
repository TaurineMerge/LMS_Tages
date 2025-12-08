/**
 * API Helper for Personal Account
 */
const API = {
    baseUrl: '/account/api/v1',
    
    getToken() {
        return localStorage.getItem('access_token');
    },
    
    setToken(token, refreshToken) {
        localStorage.setItem('access_token', token);
        if (refreshToken) {
            localStorage.setItem('refresh_token', refreshToken);
        }
    },
    
    clearTokens() {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
    },
    
    async request(endpoint, options = {}) {
        const token = this.getToken();
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };
        
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        
        const response = await fetch(`${this.baseUrl}${endpoint}`, {
            ...options,
            headers
        });
        
        if (response.status === 401) {
            const refreshed = await this.refreshToken();
            if (refreshed) {
                headers['Authorization'] = `Bearer ${this.getToken()}`;
                return fetch(`${this.baseUrl}${endpoint}`, { ...options, headers });
            } else {
                window.location.href = '/account/login';
                return;
            }
        }
        
        return response;
    },
    
    async get(endpoint) {
        const response = await this.request(endpoint);
        return response.json();
    },
    
    async post(endpoint, data) {
        const response = await this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data)
        });
        return response.json();
    },
    
    async put(endpoint, data) {
        const response = await this.request(endpoint, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
        return response.json();
    },
    
    async delete(endpoint) {
        const response = await this.request(endpoint, {
            method: 'DELETE'
        });
        if (response.status === 204) {
            return { success: true };
        }
        return response.json();
    },
    
    async refreshToken() {
        const refreshToken = localStorage.getItem('refresh_token');
        if (!refreshToken) return false;
        
        try {
            const response = await fetch(`${this.baseUrl}/auth/refresh`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh_token: refreshToken })
            });
            
            if (response.ok) {
                const data = await response.json();
                this.setToken(data.data.access_token, data.data.refresh_token);
                return true;
            }
        } catch (e) {
            console.error('Failed to refresh token:', e);
        }
        return false;
    }
};

/**
 * Toast notification helper
 */
function showToast(message, type = 'info') {
    const toastHtml = `
        <div class="alert alert-${type} alert-dismissible fade show" role="alert">
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        </div>
    `;
    
    const container = document.getElementById('toast-container');
    if (container) {
        container.innerHTML = toastHtml;
        setTimeout(() => {
            const alert = container.querySelector('.alert');
            if (alert) alert.remove();
        }, 5000);
    }
}

/**
 * Check authentication and update UI
 */
async function checkAuth() {
    const token = API.getToken();
    if (token) {
        try {
            const data = await API.get('/auth/me');
            if (data.data) {
                const user = data.data;
                const username = user.preferred_username || user.email || 'Пользователь';
                
                const navUsername = document.getElementById('nav-username');
                if (navUsername) navUsername.textContent = username;
                
                const sidebarUsername = document.querySelector('.username');
                if (sidebarUsername) sidebarUsername.textContent = username;
                
                return user;
            }
        } catch (e) {
            console.error('Auth check failed:', e);
        }
    }
    return null;
}

/**
 * Handlebars helpers
 */
if (typeof Handlebars !== 'undefined') {
    Handlebars.registerHelper('formatDate', function(dateString) {
        if (!dateString) return 'Неизвестно';
        const date = new Date(dateString);
        return date.toLocaleDateString('ru-RU', {
            day: '2-digit',
            month: 'short',
            year: 'numeric'
        });
    });
    
    Handlebars.registerHelper('formatTime', function(dateString) {
        if (!dateString) return '';
        const date = new Date(dateString);
        return date.toLocaleTimeString('ru-RU', {
            hour: '2-digit',
            minute: '2-digit'
        });
    });
    
    Handlebars.registerHelper('formatDateTime', function(dateString) {
        if (!dateString) return 'Неизвестно';
        const date = new Date(dateString);
        return date.toLocaleString('ru-RU');
    });
    
    Handlebars.registerHelper('eq', function(a, b) {
        return a === b;
    });
    
    Handlebars.registerHelper('truncate', function(str, len) {
        if (!str) return '';
        if (str.length <= len) return str;
        return str.substring(0, len) + '...';
    });
}

// Initialize on DOM ready
document.addEventListener('DOMContentLoaded', () => {
    checkAuth();
    
    // Logout handler
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', async (e) => {
            e.preventDefault();
            const refreshToken = localStorage.getItem('refresh_token');
            if (refreshToken) {
                try {
                    await API.post('/auth/logout', { refresh_token: refreshToken });
                } catch (e) {
                    console.error('Logout failed:', e);
                }
            }
            API.clearTokens();
            window.location.href = '/account/login';
        });
    }
});
