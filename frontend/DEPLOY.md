# Frontend Deployment on Render

## Quick Deploy

1. **Push to GitHub**:
```bash
git push origin main
```

2. **Create Static Site on Render**:
   - Go to https://dashboard.render.com
   - Click "New +" → "Static Site"
   - Connect your GitHub repository
   - Configure:
     - **Name**: `resume-builder-frontend`
     - **Root Directory**: `frontend`
     - **Build Command**: `npm install && npm run build`
     - **Publish Directory**: `dist`

3. **Add Environment Variable**:
   - `VITE_API_URL=https://github-resume-builder-eihv.onrender.com`

4. **Deploy**

5. **Update Backend Environment Variable**:
   - Go to your backend service on Render
   - Add: `FRONTEND_URL=https://your-frontend-url.onrender.com`
   - Redeploy backend

## Manual Deploy (Alternative)

If you prefer manual deployment:

```bash
cd frontend
npm run build
# Upload dist/ folder to any static hosting
```

## After Deployment

Update GitHub OAuth App callback URL to:
```
https://github-resume-builder-eihv.onrender.com/auth/callback
```

The frontend will be available at:
```
https://resume-builder-frontend.onrender.com
```
