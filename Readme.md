TO-DO

## Task Manager
- [x] veri tabanı tabloları oluştur
- [x] payload içeriklerini & validasyonlarını ekle
- [x] repository katmanını geliştir
- [x] repository katmanına test yaz
- [x] servis katmanını geliştir
- [x] servis katmanına test yaz
- [x] handler katmanını geliştir
- [x] server katmanını geliştir
- [x] config katmanını geliştir
- [x] loglar ekle
- [x] http error handling ekle
- [x] metric middleware i ekle
- [x] swagger dökümanı oluştur
- [x] metric endpointi sun
- [x] dockerfile oluştur
- [x] project structure uygun olsun
- [x] main.go
- [ ] schedule algoritma geliştirmesi
- [x] her bir endpointi test et
- [ ] unit testler yaz

## Cli
- [x] cobra cli altyapısı oluştur
- [x] provider alım yapısını kur
- [x] worker pool yapısı kullanılabilir
- [x] fetcher & manipulater use channels here
- [x] cobra cli tamamla
- [x] dockerfile


## Web
- [ ] Get Tasks
- [ ] Get Developers
- [ ] Get Assaignments
- [ ] Schedule Assaignments
- [ ] dockerfile

## General
- [x] docker-compose
- [ ] document
- [ ] ci
- [ ] prometheus
- [ ] name fixing
- [ ] data type checking for optimize
- [ ] migrations


docker run --rm -it --network host console:0.0.1 start  

docker run --rm --network host  \
  -e HTTP_ADDR=:8080 \
  -e HANDLER_LOG_LEVEL=debug \
  -e SERVICE_LOG_LEVEL=debug \
  -e REPOSITORY_LOG_LEVEL=debug \
  -e HTTP_SERVER_LOG_LEVEL=debug \
  -e HTTP_ALLOWED_HEADERS="*" \
  -e HTTP_ALLOWED_ORIGINS="*" \
  -e HTTP_ALLOWED_METHODS="GET,POST,PUT,DELETE,OPTIONS" \
  -e DB_HOST=localhost \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=pass \
  -e DB_NAME=task \
  -p 8080:8080 \
  task:0.0.1