import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']  
})
export class AppComponent {
  prompt: string = '';
  resultMarkdown: string = '';
  resultHtml: string = '';
  loading: boolean = false;
  error: string = '';
  activeTab: 'preview' | 'markdown' = 'preview';
  wordsJson: any = null;
  enJson: any = null;
  vocabularyHtml: string = '';

  constructor(private http: HttpClient) { }

  async sendPromptAsync(p: string): Promise<void> {
    return new Promise((resolve, reject) => {
      if (!p.trim()) {
        this.error = 'Prompt không được để trống';
        reject(new Error('Prompt không được để trống'));
        return;
      }
      this.loading = true;
      this.error = '';
      this.http.post<any>('http://localhost:8000/ask', { prompt: p })
        .subscribe({
          next: (res) => {
            this.resultMarkdown = res.markdown;
            this.resultHtml = res.html;
            this.loading = false;
            resolve();
          },
          error: (err) => {
            this.error = err.error?.error || 'Có lỗi xảy ra';
            this.loading = false;
            reject(err);
          }
        });
    });
  }

  async sendPrompt() {
    try {
      await this.sendPromptAsync(this.prompt);
      this.vocabularyHtml = this.resultHtml;
    } catch (error) {
      console.error('Lỗi gửi prompt:', error);
    }
  }

  clearResults() {
    this.resultHtml = '';
    this.resultMarkdown = '';
    this.error = '';
    this.wordsJson = null;
    this.enJson = null;
    this.vocabularyHtml = '';
    this.prompt = '';  
  }

  async pickSecretWord() {
    try {
      // Lấy nội dung hiện tại từ resultMarkdown hoặc prompt
      const currentContent = this.resultMarkdown || this.prompt;
      if (!currentContent) {
        this.error = 'Không có nội dung để phân tích';
        return;
      }

      // Bước 1: Gửi prompt để lọc từ quan trọng
      const p: string = `Từ hội thoại sau: "${currentContent}", hãy lọc ra danh sách các từ quan trọng, bỏ qua danh từ tên riêng cần học. Trả về kết quả dạng JSON mảng các từ, không cần giải thích.`;
      await this.sendPromptAsync(p);
      this.wordsJson = this.extractJsonFromMarkdown(this.resultMarkdown);

      if (this.wordsJson && Array.isArray(this.wordsJson)) {
        // Bước 2: Gửi prompt để dịch sang tiếng Anh
        const t: string = `Dịch từng từ trong danh sách sau sang tiếng Anh: ${JSON.stringify(this.wordsJson)}. Trả về JSON mảng các object có dạng {"vi": "từ tiếng Việt", "en": "từ tiếng Anh"}. Không cần giải thích.`;
        await this.sendPromptAsync(t);
        this.enJson = this.extractJsonFromMarkdown(this.resultMarkdown);
      }

      // Chuyển sang tab markdown để hiển thị kết quả
      this.activeTab = 'markdown';
    } catch (error: any) {
      console.error('Lỗi xử lý từ vựng:', error);
      this.error = 'Lỗi khi xử lý từ vựng: ' + (error.message || 'Đã xảy ra lỗi không xác định');
    }
  }

  extractJsonFromMarkdown(markdown: string): any {
    if (!markdown) return null;

    try {
      // Tìm khối JSON trong Markdown (thường được bao bởi ```json ... ```)
      const jsonMatch = markdown.match(/```json\s*([\s\S]*?)\s*```/);
      if (jsonMatch && jsonMatch[1]) {
        return JSON.parse(jsonMatch[1].trim());
      }

      // Nếu không có khối JSON, thử parse toàn bộ nội dung
      try {
        return JSON.parse(markdown.trim());
      } catch (e) {
        console.log('Không phải JSON thuần');
      }

      return null;
    } catch (e) {
      console.error('Lỗi phân tích Markdown:', e);
      return null;
    }
  }

  saveVocabulary() {
    if (!this.enJson || !Array.isArray(this.enJson)) {
      this.error = 'Không có dữ liệu từ vựng để lưu';
      return;
    }

    this.loading = true;

    const wordData = this.enJson.map(word => ({
      lang: "vi",
      content: word.vi,
      translate: word.en
    }));

    const dialogData = {
      lang: "vi",
      content: this.vocabularyHtml
    };

    const requestBody = {
      dialog: dialogData,
      words: wordData
    };

    this.http.post('http://localhost:8000/api/save', requestBody)
      .subscribe({
        next: (response) => {
          console.log('Lưu thành công:', response);
          this.loading = false;
          alert('Đã lưu từ vựng thành công!');
        },
        error: (err) => {
          console.error('Lỗi:', err);
          this.loading = false;
          this.error = 'Lỗi khi lưu từ vựng: ' + (err.error?.message || 'Đã xảy ra lỗi kết nối');
        }
      });
  }
}