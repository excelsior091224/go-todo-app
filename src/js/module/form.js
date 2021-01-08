'use strict';

// DOM Tree の構築が完了したら処理を開始します。
document.addEventListener('DOMContentLoaded', () => {
  // DOM API を利用して HTML 要素を取得します。
  const inputs = document.getElementsByTagName('input');
  const form = document.forms.namedItem('task-form');
  const saveBtn = document.querySelector('.task-form__save');
  const cancelBtn = document.querySelector('.task-form__cancel');
  const articleFormText = document.querySelector('.task-form__text');
  const articleFormTextTextArea = document.querySelector('.task-form__input--text');
  const errors = document.querySelector('.task-form__errors');
  const errorTmpl = document.querySelector('.task-form__error-tmpl').firstElementChild;
  
  // 新規作成画面か編集画面かを URL から判定します。
  const mode = { method: '', url: '' };
  if (window.location.pathname.endsWith('new')) {
    // 新規作成時の HTTP メソッドは POST を利用します。
    mode.method = 'POST';
    // 作成リクエスト、および戻るボタンの遷移先のパスは "/" になります。
    mode.url = '/tasks';
  } else if (window.location.pathname.endsWith('edit')) {
    // 更新時の HTTP メソッドは PATCH を利用します。
    mode.method = 'PATCH';
    // 更新リクエスト、および戻るボタンの遷移先のパスは "/:articleID" になります。
    mode.url = `/tasks/${window.location.pathname.split('/')[2]}`;
  }
  const { method, url } = mode;

  const csrfToken = document.getElementsByName('csrf')[0].content;

  // input 要素にフォーカスが合った状態で Enter が押されると form が送信されます。
  // 今回は Enter キーで form が送信されないように挙動を制御します。
  for (let elm of inputs) {
    elm.addEventListener('keydown', event => {
      if (event.keyCode && event.keyCode === 13) {
        // キーが押された際のデフォルトの挙動をキャンセルします。
        event.preventDefault();

        // 何もせず処理を終了します。
        return false;
      }
    });
  }

  // 前のページに戻るイベントを設定します。
  cancelBtn.addEventListener('click', event => {
    // <button> 要素クリック時のデフォルトの挙動をキャンセルします。
    event.preventDefault();

    // URL を指定して画面を遷移させます。
    window.location.href = url;
  });

  // 保存処理を実行するイベントを設定します。
  saveBtn.addEventListener('click', event => {
    event.preventDefault();

    // 前回のバリデーションエラーの表示が残っている場合は削除します。
    errors.innerHTML = null;
    
    // フォームに入力された内容を取得します。
    const fd = new FormData(form);
    
    let status;
    
    // fetch API を利用してリクエストを送信します。
    fetch(`/api${url}`, {
      method: method,
      headers: { 'X-CSRF-Token': csrfToken },
      body: fd
    })
    .then(res => {
      status = res.status;
      return res.json();
    })
    .then(body => {
      console.log(JSON.stringify(body));
      
      if (status === 200) {
        // 成功時は一覧画面に遷移させます。
        window.location.href = url;
      }
  
      if (body.ValidationErrors) {
        // バリデーションエラーがある場合の処理をここに記載します。
        showErrors(body.ValidationErrors);
      }
    })
    .catch(err => console.error(err));
  });

  // バリデーションエラーを表示する関数
  const showErrors = messages => {
    // 引数の値が配列であることを確認します。
    if (Array.isArray(messages) && messages.length != 0) {
      // 複数メッセージを格納するためのフラグメントを作成します。
      const fragment = document.createDocumentFragment();

      // メッセージをループ処理します。
      messages.forEach(message => {
        // 単一メッセージを格納するためのフラグメントを作成します。
        const frag = document.createDocumentFragment();
  
        // テンプレートをクローンしてフラグメントに追加します。
        frag.appendChild(errorTmpl.cloneNode(true));

        // エラー要素にメッセージをセットします。
        frag.querySelector('.task-form__error').innerHTML = message;

        // エラー要素を親フラグメントに追加します。
        fragment.appendChild(frag);
      });

      // エラーメッセージの表示エリア（要素）にメッセージを追加します。
      errors.appendChild(fragment);
    }
  };
});