import React, { useEffect, useState } from 'react';
import './App.css';

interface HealthResponse {
  status: string;
  service: string;
}

function App() {
  const [healthStatus, setHealthStatus] = useState<HealthResponse | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // バックエンドヘルスチェック
    fetch('/api/health')
      .then(response => response.json())
      .then((data: HealthResponse) => {
        setHealthStatus(data);
        setLoading(false);
      })
      .catch(error => {
        console.error('API接続エラー:', error);
        setLoading(false);
      });
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* ヘッダー */}
      <header className="bg-blue-600 text-white shadow-lg">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold">建築業統合管理システム</h1>
            </div>
            <div className="flex items-center space-x-4">
              <div className="text-sm">
                {loading ? (
                  <span className="text-blue-200">接続確認中...</span>
                ) : healthStatus ? (
                  <span className="text-green-200">✓ API接続正常</span>
                ) : (
                  <span className="text-red-200">✗ API接続エラー</span>
                )}
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* メインコンテンツ */}
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="border-4 border-dashed border-gray-200 rounded-lg h-96">
            <div className="flex items-center justify-center h-full">
              <div className="text-center">
                <h2 className="text-2xl font-bold text-gray-900 mb-4">
                  開発環境構築完了！
                </h2>
                <p className="text-gray-600 mb-6">
                  建築業向け統合管理システムの開発をスタートしましょう
                </p>

                {/* システム状況 */}
                <div className="bg-white p-6 rounded-lg shadow max-w-md mx-auto">
                  <h3 className="text-lg font-semibold mb-4">システム状況</h3>
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span>フロントエンド:</span>
                      <span className="text-green-600">✓ 起動中</span>
                    </div>
                    <div className="flex justify-between">
                      <span>バックエンドAPI:</span>
                      {healthStatus ? (
                        <span className="text-green-600">✓ 接続正常</span>
                      ) : (
                        <span className="text-red-600">✗ 接続エラー</span>
                      )}
                    </div>
                    <div className="flex justify-between">
                      <span>データベース:</span>
                      <span className="text-green-600">✓ 初期化完了</span>
                    </div>
                  </div>
                </div>

                {/* 次のステップ */}
                <div className="mt-8">
                  <h3 className="text-lg font-semibold mb-4">次のステップ</h3>
                  <div className="grid grid-cols-1 gap-4 max-w-lg mx-auto">
                    <button className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition-colors">
                      プロジェクト管理画面を作成
                    </button>
                    <button className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 transition-colors">
                      資材管理機能を実装
                    </button>
                    <button className="bg-purple-600 text-white px-4 py-2 rounded hover:bg-purple-700 transition-colors">
                      ユーザー認証を追加
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

export default App;
