document.addEventListener('DOMContentLoaded', async () => {
    // Проверка аутентификации
    const token = localStorage.getItem('jwt');
    // if (!token) {
    //     window.location.href = '/frontend/login.html';
    //     return;
    // }

    // Элементы UI
    const logoutBtn = document.getElementById('logout-btn');
    const updateAvatarBtn = document.getElementById('update-avatar-btn');
    const updateBioBtn = document.getElementById('update-bio-btn');

    // Загрузка данных профиля
    async function loadProfile() {
        try {
            const response = await fetch('/api/profile', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.status === 401) {
                // Неавторизован
                localStorage.removeItem('jwt');
                window.location.href = '/frontend/login.html';
                return;
            }

            if (!response.ok) {
                throw new Error('Failed to load profile');
            }

            const data = await response.json();
            displayProfile(data.profile);
        } catch (error) {
            console.error('Error:', error);
            alert('Failed to load profile');
        }
    }

    function displayProfile(profile) {
        document.getElementById('username').textContent = profile.username;
        document.getElementById('email').textContent = profile.email;
        document.getElementById('joined-at').textContent = new Date(profile.created_at).toLocaleDateString();
        document.getElementById('bio-input').value = profile.bio || '';
        
        if (profile.avatar_url) {
            document.getElementById('avatar-img').src = profile.avatar_url;
        }
        
        // Статистика (если есть)
        if (profile.stats) {
            document.getElementById('wins').textContent = profile.stats.wins;
            document.getElementById('losses').textContent = profile.stats.losses;
        }
    }

    // Обновление аватара
    // updateAvatarBtn.addEventListener('click', async () => {
    //     const newAvatarUrl = document.getElementById('avatar-url').value.trim();
    //     if (!newAvatarUrl) return;

    //     try {
    //         const response = await fetch('/api/profile', {
    //             method: 'PUT',
    //             headers: {
    //                 'Content-Type': 'application/json',
    //                 'Authorization': `Bearer ${token}`
    //             },
    //             body: JSON.stringify({
    //                 avatar_url: newAvatarUrl
    //             })
    //         });

    //         if (!response.ok) {
    //             throw new Error('Failed to update avatar');
    //         }

    //         document.getElementById('avatar-img').src = newAvatarUrl;
    //         document.getElementById('avatar-url').value = '';
    //         showToast('Avatar updated successfully!');
    //     } catch (error) {
    //         console.error('Error:', error);
    //         showToast(error.message, 'error');
    //     }
    // });

    // Обновление био
    updateBioBtn.addEventListener('click', async () => {
        const newBio = document.getElementById('bio-input').value;

        try {
            const response = await fetch('/api/profile', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    bio: newBio
                })
            });

            if (!response.ok) {
                throw new Error('Failed to update bio');
            }

            showToast('Journal updated successfully!');
        } catch (error) {
            console.error('Error:', error);
            showToast(error.message, 'error');
        }
    });

    // Выход
    logoutBtn.addEventListener('click', () => {
        localStorage.removeItem('jwt');
        localStorage.removeItem('user');
        window.location.href = '/frontend/login.html';
    });

    // Вспомогательные функции
    function showToast(message, type = 'success') {
        const toast = document.createElement('div');
        toast.className = `toast toast-${type}`;
        toast.textContent = message;
        document.body.appendChild(toast);
        
        setTimeout(() => {
            toast.remove();
        }, 3000);
    }

    // Инициализация
    loadProfile();
});