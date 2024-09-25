var gulp = require('gulp');
var watch = require('gulp-watch');
var cleanCSS = require('gulp-clean-css');
var uglify = require('gulp-uglify-es').default;
var rename = require('gulp-rename');
var concat = require('gulp-concat');
var merge = require('merge-stream');
var scss = require('gulp-sass');
var fs = require('fs');

function version(){
  var now = new Date(),
    Y = now.getFullYear(),
    m = now.getMonth()+1,
    d = now.getDate(),
    H = now.getHours(),
    i = now.getMinutes(),
    s = now.getSeconds();

    if(H < 10) {
        H = '0' + H;
    }

    if(i < 10) {
        i = '0' + i;
    }

    if(s < 10) {
        s = '0' + s;
    }

    return String(10000*Y + 100*m + d + '.' + H + i + s);
}

function setVersion(type) {
  fs.writeFile(`./assets/bust/${type}.txt`, version(), function(err) {
    if(err) {
        return console.log(err);
    }
  }); 
}

function buildMainCSS(){

  var full = gulp.src([
    'assets/src/scss/main.scss',
    'assets/src/scss/navbar.scss'
  ])
  . pipe(scss())
  . pipe(concat('main.css'))
  . pipe(gulp.dest('static/dist/css'));

  var min = gulp.src([
    'assets/src/scss/main.scss',
    'assets/src/scss/navbar.scss'
  ])
  . pipe(scss())
  . pipe(cleanCSS())
  . pipe(concat('main.min.css'))
  . pipe(gulp.dest('static/dist/css'));

  setVersion('css');
  return merge(min, full);
}

function buildHomeCSS(){

  var full = gulp.src([
    'assets/src/scss/home-page.scss',
    'assets/src/scss/navbar.scss'
  ])
  . pipe(scss())
  . pipe(concat('home-page.css'))
  . pipe(gulp.dest('static/dist/css'));

  var min = gulp.src([
    'assets/src/scss/home-page.scss',
    'assets/src/scss/navbar.scss'
  ])
  . pipe(scss())
  . pipe(cleanCSS())
  . pipe(concat('home-page.min.css'))
  . pipe(gulp.dest('static/dist/css'));

  setVersion('css');
  return merge(full, min);
}


function buildAdminCSS(){

  var full = gulp.src([
    'assets/src/scss/admin.scss'
  ])
  . pipe(scss())
  . pipe(concat('admin.min.css'))
  . pipe(gulp.dest('static/dist/css'));

  var min = gulp.src([
    'assets/src/scss/admin.scss'
  ])
  . pipe(scss())
  . pipe(cleanCSS())
  . pipe(concat('admin.min.css'))
  . pipe(gulp.dest('static/dist/css'));

  setVersion('css');
  return merge(full, min);
}


function buildHomeJS() {

  var full = gulp.src([
    'assets/src/js/main.js',
    'assets/src/js/home.js',
  ])
  .pipe(concat('home-page.js'))
  .pipe(gulp.dest('static/dist/js'));

  var min = gulp.src([
    'assets/src/js/main.js',
    'assets/src/js/home.js',
  ])
  .pipe(concat('home-page.min.js'))
  .pipe(uglify())
  .pipe(gulp.dest('static/dist/js'));

  setVersion('css');
  return merge(full, min);
}

function buildArticleJS() {

  var full = gulp.src([
    'assets/src/js/main.js',
    'assets/src/js/article.js',
  ])
  .pipe(concat('article.js'))
  .pipe(gulp.dest('static/dist/js'));

  var min = gulp.src([
    'assets/src/js/main.js',
    'assets/src/js/article.js',
  ])
  .pipe(concat('article.min.js'))
  .pipe(uglify())
  .pipe(gulp.dest('static/dist/js'));

  setVersion('js');
  return merge(min, full);

}

function buildHashJS(){

  var hash = gulp.src([
    'node_modules/jssha/src/sha.js',
    'node_modules/js-md5/src/md5.js',
    'assets/src/js/tools/hash.js'
  ])
  .pipe(concat('hash.min.js'))
  .pipe(uglify())
  .pipe(gulp.dest('static/dist/js/tools'));

  setVersion('js');
  return merge(hash);
}

function buildAdminJS() {

  var full = gulp.src([
    'assets/src/js/admin.js'
  ])
  .pipe(concat('admin.js'))
  .pipe(uglify())
  .pipe(gulp.dest('static/dist/js'));

  var min = gulp.src([
    'assets/src/js/admin.js'
  ])
  .pipe(concat('admin.min.js'))
  .pipe(uglify())
  .pipe(gulp.dest('static/dist/js'));

  setVersion('js');
  return merge(full, min);
}


function watchFiles() {
  gulp.watch(['./assets/src/scss/main.scss', './assets/src/scss/navbar.scss'], buildMainCSS);
  gulp.watch(['./assets/src/scss/main.scss', './assets/src/scss/navbar.scss', './assets/src/scss/home-page.scss'], buildHomeCSS);
  gulp.watch(['./assets/src/scss/admin.scss', './assets/src/scss/navbar.scss'], buildAdminCSS);
  gulp.watch(['./assets/src/js/main.js', './assets/src/js/home.js'], buildHomeJS);
  gulp.watch(['./assets/src/js/main.js', './assets/src/js/article.js'], buildArticleJS);
  gulp.watch('./assets/src/js/admin.js', buildAdminJS);
}

gulp.task('build-admin-css', buildAdminCSS);

gulp.task('build-main-css', buildMainCSS); 

gulp.task('build-home-css', buildHomeCSS); 

gulp.task('build-home-js', buildHomeJS);

gulp.task('build-article-js', buildArticleJS);

gulp.task('build-admin-js', buildAdminJS);

gulp.task('build-hash-js', buildHashJS);

gulp.task('default', watchFiles);

gulp.task('watch', watchFiles);
