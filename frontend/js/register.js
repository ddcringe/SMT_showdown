document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('register-form');
    const submitBtn = document.getElementById('submit-btn');
    const spinner = document.getElementById('spinner');

    // Валидация в реальном времени
    document.getElementById('password').addEventListener('input', validatePassword);
    document.getElementById('confirm-password').addEventListener('input', validatePasswordMatch);

    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        if (!validateForm()) {
            return;
        }

        const formData = {
            username: document.getElementById('username').value.trim(),
            email: document.getElementById('email').value.trim(),
            password: document.getElementById('password').value
        };

        try {
            // Показываем индикатор загрузки
            submitBtn.disabled = true;
            spinner.classList.remove('hidden');
            submitBtn.querySelector('.button-text').textContent = 'Processing...';

            const response = await fetch('/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            });

            const data = await response.json();

            if (!response.ok) {
                handleApiErrors(data);
                return;
            }

            // Успешная регистрация
            showSuccessMessage('Registration successful! Redirecting to login...');
            
            // Сохраняем токен, если бекенд сразу его возвращает
            if (data.token) {
                localStorage.setItem('jwt', data.token);
            }
            
            // Перенаправляем через 1.5 секунды
            setTimeout(() => {
                window.location.href = '/frontend/login.html';
            }, 1500);
            
        } catch (error) {
            console.error('Registration error:', error);
            showErrorMessage('Network error. Please try again later.');
        } finally {
            submitBtn.disabled = false;
            spinner.classList.add('hidden');
            submitBtn.querySelector('.button-text').textContent = 'Create Account';
        }
    });

    function validateForm() {
        let isValid = true;
        
        // Очищаем предыдущие ошибки
        clearErrors();
        
        // Проверка имени пользователя
        const username = document.getElementById('username').value.trim();
        if (!username || username.length < 3) {
            showError('username-error', 'Username must be at least 3 characters');
            isValid = false;
        }
        
        // Проверка email
        const email = document.getElementById('email').value.trim();
        if (!email || !validateEmail(email)) {
            showError('email-error', 'Please enter a valid email');
            isValid = false;
        }
        
        // Проверка пароля
        const password = document.getElementById('password').value;
        if (!password || password.length < 8) {
            showError('password-error', 'Password must be at least 8 characters');
            isValid = false;
        }
        
        // Проверка совпадения паролей
        const confirmPassword = document.getElementById('confirm-password').value;
        if (password !== confirmPassword) {
            showError('confirm-error', 'Passwords do not match');
            isValid = false;
        }
        
        return isValid;
    }

    function validatePassword() {
        const password = document.getElementById('password').value;
        const errorElement = document.getElementById('password-error');
        
        if (password.length > 0 && password.length < 8) {
            showError('password-error', 'Password too short');
        } else {
            clearError('password-error');
        }
    }

    function validatePasswordMatch() {
        const password = document.getElementById('password').value;
        const confirmPassword = document.getElementById('confirm-password').value;
        
        if (confirmPassword.length > 0 && password !== confirmPassword) {
            showError('confirm-error', 'Passwords do not match');
        } else {
            clearError('confirm-error');
        }
    }

    function validateEmail(email) {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    }

    function handleApiErrors(data) {
        if (data.errors) {
            // Обработка ошибок валидации с сервера
            Object.keys(data.errors).forEach(field => {
                showError(`${field}-error`, data.errors[field]);
            });
        } else {
            showErrorMessage(data.error || 'Registration failed. Please try again.');
        }
    }

    function showError(elementId, message) {
        const element = document.getElementById(elementId);
        if (element) {
            element.textContent = message;
            element.style.display = 'block';
        }
    }

    function clearError(elementId) {
        const element = document.getElementById(elementId);
        if (element) {
            element.textContent = '';
            element.style.display = 'none';
        }
    }

    function clearErrors() {
        const errorElements = document.querySelectorAll('.error-message');
        errorElements.forEach(el => {
            el.textContent = '';
            el.style.display = 'none';
        });
    }

    function showSuccessMessage(message) {
        alert(message); // Можно заменить на красивый toast-уведомление
    }

    function showErrorMessage(message) {
        alert(message); // Можно заменить на красивый toast-уведомление
    }
});