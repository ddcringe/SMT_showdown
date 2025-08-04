document.addEventListener('DOMContentLoaded', () => {
    loadProfile();
});

async function loadProfile() {
    try {
        const token = localStorage.getItem('jwt');
        if (!token) {
            window.location.href = '/login.html';
            return;
        }

        const response = await fetch('/api/profile', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to load profile');
        }

        const data = await response.json();
        displayProfile(data.profile);
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to load profile: ' + error.message);
    }
}

function displayProfile(profile) {
    document.getElementById('username').textContent = profile.username;
    document.getElementById('email').textContent = profile.email;
    document.getElementById('joined-at').textContent = new Date(profile.created_at).toLocaleDateString();
    
    const bioInput = document.getElementById('bio-input');
    bioInput.value = profile.bio || '';
    
    const avatarImg = document.getElementById('avatar-img');
    avatarImg.src = profile.avatar_url || '/frontend/images/default-avatar.jpg';
}

async function updateBio() {
    const newBio = document.getElementById('bio-input').value;
    const token = localStorage.getItem('jwt');

    try {
        const response = await fetch('/api/profile', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ bio: newBio })
        });

        if (!response.ok) {
            throw new Error('Failed to update bio');
        }

        alert('Bio updated successfully!');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to update bio: ' + error.message);
    }
}

async function updateAvatar() {
    const newAvatarUrl = document.getElementById('avatar-url').value;
    const token = localStorage.getItem('jwt');

    try {
        const response = await fetch('/api/profile', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ avatar_url: newAvatarUrl })
        });

        if (!response.ok) {
            throw new Error('Failed to update avatar');
        }

        document.getElementById('avatar-img').src = newAvatarUrl;
        alert('Avatar updated successfully!');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to update avatar: ' + error.message);
    }
}