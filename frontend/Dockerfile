# 1. Node.js tabanlı bir image kullanarak uygulamayı build ediyoruz
FROM node:16 AS build

# 2. Çalışma dizinini ayarla
WORKDIR /app

# 3. Package.json ve package-lock.json dosyalarını kopyala
COPY package*.json ./

# 4. Gerekli bağımlılıkları yükle
RUN npm install

# 5. Uygulama kaynak kodunu kopyala
COPY . .

# 6. React uygulamasını build et
RUN npm run build

# 7. Nginx tabanlı bir image kullanarak build edilen dosyaları serve ediyoruz
FROM nginx:alpine

# 8. Build edilen dosyaları Nginx'in varsayılan dizinine kopyala
COPY --from=build /app/build /usr/share/nginx/html

# 9. Nginx için portu expose et
EXPOSE 80

# 10. Nginx'i başlat
CMD ["nginx", "-g", "daemon off;"]