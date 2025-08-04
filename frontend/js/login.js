document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('login-form');
    const loginBtn = document.getElementById('login-btn');
    const spinner = document.getElementById('spinner');

    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const email = document.getElementById('email').value.trim();
        const password = document.getElementById('password').value;

        try {
            // Показываем индикатор загрузки
            loginBtn.disabled = true;
            spinner.classList.remove('hidden');
            loginBtn.querySelector('.button-text').textContent = 'Authenticating...';

            const response = await fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    email,
                    password
                })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Login failed');
            }

            // Сохраняем токен
            localStorage.setItem('jwt', data.token);
            localStorage.setItem('user', JSON.stringify({
                id: data.user_id,
                email: email
            }));

            // Перенаправляем на профиль
            window.location.href = '/frontend/profile.html';

        } catch (error) {
            console.error('Login error:', error);
            document.getElementById('password-error').textContent = error.message;
            document.getElementById('password-error').style.display = 'block';
        } finally {
            loginBtn.disabled = false;
            spinner.classList.add('hidden');
            loginBtn.querySelector('.button-text').textContent = 'Log In';
        }
    });
});