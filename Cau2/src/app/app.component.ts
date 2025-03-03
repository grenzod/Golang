import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  voices = [
    'en-US-AndrewMultilingualNeural',
    'en-US-ChristopherNeural',
    'en-US-EricNeural',
    'vi-VN-HoaiMyNeural',
    'vi-VN-NamMinhNeural'
  ];
  text: string = '';
  showResult: boolean = false;

  selectedVoiceA: string = this.voices[0];
  selectedVoiceB: string = this.voices[3];

  conversation = [
    { speaker: 'A', text: 'Chào Lan! Mình là James, đến từ Hoa Kỳ. Rất vui được gặp bạn.' },
    { speaker: 'B', text: 'Chào James! Mình là Lan, đến từ Việt Nam. Rất vui được làm quen với bạn.' },
    { speaker: 'A', text: 'Bạn làm nghề gì vậy, Lan?' },
    { speaker: 'B', text: 'Mình là cô giáo dạy ngoại ngữ. Còn bạn?' },
    { speaker: 'A', text: 'Mình là kỹ sư hàng không.' },
    { speaker: 'B', text: 'Nghe thú vị quá! Bạn đến Việt Nam lâu chưa?' },
    { speaker: 'A', text: 'Mình mới đến đây được vài ngày.' },
    { speaker: 'B', text: 'Hy vọng bạn sẽ thích Việt Nam!' },
    { speaker: 'A', text: 'Cảm ơn Lan!' }
  ];

  generateSSML() {
    let ssml = '<speak xml:lang="vi-VN">';
    this.conversation.forEach(line => {
      const voice = line.speaker === 'A' ? this.selectedVoiceA : this.selectedVoiceB;
      ssml += `<voice name="${voice}">${line.text}</voice>`;
    });
    ssml += '</speak>';
    this.text = ssml;
    this.showResult = true;
  }

  downloadSSML() {
    const blob = new Blob([this.text], { type: 'text/xml' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'conversation.ssml';
    a.click();
    window.URL.revokeObjectURL(url);
  }
}