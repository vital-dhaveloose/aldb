import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { NavPaneComponent } from "./components/nav-pane/nav-pane.component";
import { ActivityVersion } from './models/activityVersion';
import { AldbService } from './services/aldb.service';

@Component({
    selector: 'app-root',
    standalone: true,
    templateUrl: './app.component.html',
    styleUrl: './app.component.css',
    imports: [CommonModule, RouterOutlet, NavPaneComponent]
})
export class AppComponent {
  title = 'canviga-web';
}
