import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  prompt: string = '';
  resultMarkdown: string = '';
  resultHtml: string = '';
  loading: boolean = false;
  error: string = '';
  activeTab: 'preview' | 'markdown' = 'preview';

  constructor(private http: HttpClient) {}

  sendPrompt() {
    if (!this.prompt.trim()) {
      this.error = 'Prompt không được để trống';
      return;
    }
    this.loading = true;
    this.error = '';

    // Gọi API backend tại endpoint /ask
    this.http.post<any>('http://localhost:8000/ask', { prompt: this.prompt })
      .subscribe(
        res => {
          this.resultMarkdown = res.markdown;
          this.resultHtml = res.html;
          this.loading = false;
        },
        err => {
          this.error = err.error.error || 'Có lỗi xảy ra';
          this.loading = false;
        }
      );
  }

    // Thêm phương thức này
    clearResults() {
      this.resultHtml = '';
      this.resultMarkdown = '';
      this.error = '';
    }
    
    // Thêm phương thức sao chép vào clipboard
    copyToClipboard(type: 'html' | 'markdown') {
      const textToCopy = type === 'html' ? this.resultHtml : this.resultMarkdown;
      
      navigator.clipboard.writeText(textToCopy).then(() => {
        // Có thể thêm thông báo đã sao chép thành công
        console.log('Đã sao chép vào clipboard');
      }, (err) => {
        console.error('Không thể sao chép: ', err);
      });
    }
}
